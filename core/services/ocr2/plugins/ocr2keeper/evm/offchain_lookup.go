package evm

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2keepers/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
)

type MercuryLookup struct {
	feedLabel  string
	feeds      []string
	queryLabel string
	query      *big.Int
	extraData  []byte
}

// MercuryMultiResponse TODO guessing on this
type MercuryMultiResponse struct {
	ChainlinkBlobs []string `json:"chainlinkBlob"`
}

type MercuryResponse struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

type MercuryBytes struct {
	Index int
	Bytes []byte
}

// mercuryLookup looks through check upkeep results looking for any that need off chain lookup
func (r *EvmRegistry) mercuryLookup(ctx context.Context, upkeepResults []types.UpkeepResult) ([]types.UpkeepResult, error) {
	for i := range upkeepResults {
		// if its another reason continue/skip
		if upkeepResults[i].FailureReason != UPKEEP_FAILURE_REASON_TARGET_CHECK_REVERTED {
			continue
		}

		block, upkeepId, err := blockAndIdFromKey(upkeepResults[i].Key)
		if err != nil {
			r.lggr.Error("[MercuryLookup] error getting block and upkeep id:", err)
			continue
		}

		// checking if this upkeep is in cooldown from api errors
		_, onIce := r.mercury.cooldownCache.Get(upkeepId.String())
		if onIce {
			r.lggr.Infof("[MercuryLookup] cooldown Skipping UpkeepId: %s\n", upkeepId)
			continue
		}

		// if it doesn't decode to the offchain custom error continue/skip
		mercuryLookup, err := r.decodeMercuryLookup(upkeepResults[i].PerformData)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR
			r.lggr.Debug("[MercuryLookup] not an offchain revert decodeMercuryLookup:", err)
			continue
		}
		r.lggr.Debugf("[MercuryLookup]: %+v\n", mercuryLookup)

		opts, err := r.buildCallOpts(ctx, block)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR
			r.lggr.Error("[MercuryLookup] buildCallOpts:", err)
			continue
		}
		// need upkeep info for offchainConfig and to hit callback
		upkeepInfo, err := r.getUpkeepInfo(upkeepId, opts)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR
			r.lggr.Error("[MercuryLookup] GetUpkeep:", err)
			continue
		}

		// 	do the http request
		values, err := r.doRequest(mercuryLookup, upkeepId)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR
			r.lggr.Error("[MercuryLookup] doRequest:", err)
			continue
		}

		needed, performData, err := r.mercuryLookupCallback(ctx, mercuryLookup, values, upkeepInfo, opts)
		if err != nil {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_MERCURY_LOOKUP_ERROR
			r.lggr.Error("[MercuryLookup] mercuryLookupCallback=", err)
			continue
		}
		if !needed {
			upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_UPKEEP_NOT_NEEDED
			r.lggr.Debug("[MercuryLookup] callback reports upkeep not needed")
			continue
		}

		// success!
		upkeepResults[i].FailureReason = UPKEEP_FAILURE_REASON_NONE
		upkeepResults[i].State = types.Eligible
		upkeepResults[i].PerformData = performData
		r.lggr.Debugf("[MercuryLookup] Success: %+v\n", upkeepResults[i])
	}
	return upkeepResults, nil
}

func (r *EvmRegistry) getUpkeepInfo(upkeepId *big.Int, opts *bind.CallOpts) (keeper_registry_wrapper2_0.UpkeepInfo, error) {
	zero := common.Address{}
	var err error
	var upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo
	u, found := r.mercury.upkeepCache.Get(upkeepId.String())
	if found {
		upkeepInfo = u.(keeper_registry_wrapper2_0.UpkeepInfo)
		r.lggr.Debugf("[MercuryLookup] cache hit UpkeepInfo: %+v\n", upkeepInfo)
	} else {
		upkeepInfo, err = r.registry.GetUpkeep(opts, upkeepId)
		if err != nil {
			return upkeepInfo, err
		}
		if upkeepInfo.Target == zero {
			return upkeepInfo, errors.New("upkeepInfo should not be nil")
		}
		r.lggr.Debugf("[MercuryLookup] cache miss UpkeepInfo: %+v\n", upkeepInfo)
		r.mercury.upkeepCache.Set(upkeepId.String(), upkeepInfo, cache.DefaultExpiration)
	}
	return upkeepInfo, nil
}

// decodeMercuryLookup decodes the revert error ChainlinkAPIFetch(string query, bytes extraData, string[] jsonFields, bytes4 callbackSelector)
func (r *EvmRegistry) decodeMercuryLookup(data []byte) (MercuryLookup, error) {
	e := r.mercury.abi.Errors["MercuryLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return MercuryLookup{}, errors.Wrapf(err, "unpack error")
	}
	errorParameters := unpack.([]interface{})

	return MercuryLookup{
		feedLabel:  *abi.ConvertType(errorParameters[0], new(string)).(*string),
		feeds:      *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		queryLabel: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		query:      *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		extraData:  *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

// mercuryLookupCallback calls the callback(string[] memory chainlinkBlobHex, bytes memory extraData) specified by the
// 4-byte selector from the revert. the return will match check telling us if the upkeep is needed and what the perform data is
func (r *EvmRegistry) mercuryLookupCallback(ctx context.Context, mercuryLookup MercuryLookup, values [][]byte, upkeepInfo keeper_registry_wrapper2_0.UpkeepInfo, opts *bind.CallOpts) (bool, []byte, error) {
	payload, err := r.mercury.abi.Pack("mercuryCallback", values, mercuryLookup.extraData)
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback args pack error")
	}

	checkUpkeepGasLimit := uint32(200000) + uint32(6500000) + uint32(300000) + upkeepInfo.ExecuteGas
	callbackMsg := ethereum.CallMsg{
		From: r.addr,             // registry addr
		To:   &upkeepInfo.Target, // upkeep addr
		Gas:  uint64(checkUpkeepGasLimit),
		Data: payload,
	}

	callbackResp, err := r.client.CallContract(ctx, callbackMsg, opts.BlockNumber)
	if err != nil {
		return false, nil, errors.Wrapf(err, "call contract callback error")
	}

	typBytes, err := abi.NewType("bytes", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new bytes type error")
	}
	boolTyp, err := abi.NewType("bool", "", nil)
	if err != nil {
		return false, nil, errors.Wrapf(err, "abi new bool type error")
	}
	callbackOutput := abi.Arguments{
		{Name: "upkeepNeeded", Type: boolTyp},
		{Name: "performData", Type: typBytes},
	}
	unpack, err := callbackOutput.Unpack(callbackResp)
	if err != nil {
		return false, nil, errors.Wrapf(err, "callback output unpack error")
	}

	upkeepNeeded := *abi.ConvertType(unpack[0], new(bool)).(*bool)
	if !upkeepNeeded {
		return false, nil, nil
	}
	performData := *abi.ConvertType(unpack[1], new([]byte)).(*[]byte)
	return true, performData, nil
}

func (r *EvmRegistry) doRequest(mercuryLookup MercuryLookup, upkeepId *big.Int) ([][]byte, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// TODO when mercury has multi feed endpoint. we can use this instead of below
	//multiFeed, err := r.multiFeedRequest(&client, upkeepId, mercuryLookup)
	//if err != nil {
	//	return nil, err
	//}
	//return multiFeed, nil

	// TODO this is if Mercury doesn't have the multi feeds endpoint
	ch := make(chan MercuryBytes)
	for i := range mercuryLookup.feeds {
		// if mercury ends up providing an endpoint to do all feeds at once great, if not...
		go r.singleFeedRequest(&client, ch, upkeepId, i, mercuryLookup)
	}
	results := make([][]byte, len(mercuryLookup.feeds))
	for i := 0; i < len(results); i++ {
		m := <-ch
		results[m.Index] = m.Bytes
	}

	return results, nil
}

func (r *EvmRegistry) singleFeedRequest(client *http.Client, ch chan<- MercuryBytes, upkeepId *big.Int, index int, mercuryLookup MercuryLookup) {
	req, err := http.NewRequest("GET", r.mercury.url, nil)
	if err != nil {
		ch <- MercuryBytes{Index: index}
		return
	}
	q := url.Values{}
	q.Add(mercuryLookup.feedLabel, mercuryLookup.feeds[index])
	q.Add(mercuryLookup.queryLabel, mercuryLookup.query.String())
	q.Add("upkeepID", upkeepId.String())
	req.URL.RawQuery = q.Encode()

	signature := generateHMAC("GET", req.URL.String(), nil, r.mercury.clientID, r.mercury.clientKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.clientID)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	// TODO no test client ID/Key yet. Just return correct json for testing
	//resp, err := client.Do(req)
	//if err != nil {
	//	r.setCachesOnAPIErr(upkeepId)
	//	ch <- MercuryBytes{Index: index}
	//	return
	//}
	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	r.setCachesOnAPIErr(upkeepId)
	//	ch <- MercuryBytes{Index: index}
	//	return
	//}
	//// if we get a 403 permission issue we can put them on a longer cooldown to avoid spamming mercury
	//// if http response code is 4xx/5xx then put in cool down
	//if resp.StatusCode >= 400 {
	//	r.setCachesOnAPIErr(upkeepId)
	//}
	//var m MercuryResponse
	//err = json.Unmarshal(body, &m)
	//if err != nil {
	//	ch <- MercuryBytes{Index: index}
	//	return
	//}

	m := MercuryResponse{ChainlinkBlob: "0x000189dbcc9287f900f77bea62d479cfd70ec8073692ca911fe306cd5bcf8d6d0000000000000000000000000000000000000000000000000000000000100e58c41df85f0fb47f78779e68b0a0dbefe8d626446b286aebf96741ae274cf3a49d00000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001e0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800fb2e5752573270cb04af9a1ebafc82b67f09a7408217b3b3cb81fe24eb0912900000000000000000000000000000000000000000000000000000000637ce26100000000000000000000000000000000000000000000000000000000002bbecc00000000000000000000000000000000000000000000003cb243ded70e1d0000000000000000000000000000000000000000000000000000000000000000000243d68b3eda5fb3a526daffdab2bf9978f89a6c1a5de19020fd047b80692b67a2712668bc873498a2a69f38ea7874f7e3511baa0637af1b1b304e4cae2bf87c1400000000000000000000000000000000000000000000000000000000000000021e1e506899f5ea70c67ea458042e08a9b45f888f67fa31ae286d760e1c967d14210469b4efc32630d4af00842811b739d7816439abbc9ee48f2d463f3f657fd7"}
	blobBytes, err := hexutil.Decode(m.ChainlinkBlob)
	if err != nil {
		ch <- MercuryBytes{Index: index}
		return
	}
	ch <- MercuryBytes{Index: index, Bytes: blobBytes}
	return
}

func (r *EvmRegistry) multiFeedRequest(client *http.Client, upkeepId *big.Int, mercuryLookup MercuryLookup) ([][]byte, error) {
	req, err := http.NewRequest("GET", r.mercury.url, nil)
	if err != nil {
		return [][]byte{}, err
	}
	q := url.Values{}
	feeds := strings.Join(mercuryLookup.feeds, ",")
	q.Add(mercuryLookup.feedLabel, feeds)
	q.Add(mercuryLookup.queryLabel, mercuryLookup.query.String())
	q.Add("upkeepID", upkeepId.String())
	req.URL.RawQuery = q.Encode()

	signature := generateHMAC("GET", req.URL.String(), nil, r.mercury.clientID, r.mercury.clientKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.mercury.clientID)
	req.Header.Set("X-Authorization-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("X-Authorization-Signature-SHA256", signature)

	// TODO no test client ID/Key yet. Just return correct json for testing
	resp, err := client.Do(req)
	if err != nil {
		r.setCachesOnAPIErr(upkeepId)
		return [][]byte{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.setCachesOnAPIErr(upkeepId)
		return [][]byte{}, err
	}
	// TODO ? if we get a 403 permission issue we can put them on a longer cooldown to avoid spamming mercury
	// if http response code is 4xx/5xx then put in cool down
	if resp.StatusCode >= 400 {
		r.setCachesOnAPIErr(upkeepId)
	}
	var m MercuryMultiResponse
	err = json.Unmarshal(body, &m)
	if err != nil {
		return [][]byte{}, err
	}

	// TODO this is just a guess of what they will return
	//m := MercuryMultiResponse{
	//	ChainlinkBlobs: []string{
	//		"0x000189dbcc9287f900f77bea62d479cfd70ec8073692ca911fe306cd5bcf8d6d0000000000000000000000000000000000000000000000000000000000100e58c41df85f0fb47f78779e68b0a0dbefe8d626446b286aebf96741ae274cf3a49d00000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001e0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800fb2e5752573270cb04af9a1ebafc82b67f09a7408217b3b3cb81fe24eb0912900000000000000000000000000000000000000000000000000000000637ce26100000000000000000000000000000000000000000000000000000000002bbecc00000000000000000000000000000000000000000000003cb243ded70e1d0000000000000000000000000000000000000000000000000000000000000000000243d68b3eda5fb3a526daffdab2bf9978f89a6c1a5de19020fd047b80692b67a2712668bc873498a2a69f38ea7874f7e3511baa0637af1b1b304e4cae2bf87c1400000000000000000000000000000000000000000000000000000000000000021e1e506899f5ea70c67ea458042e08a9b45f888f67fa31ae286d760e1c967d14210469b4efc32630d4af00842811b739d7816439abbc9ee48f2d463f3f657fd7",
	//		"0x000189dbcc9287f900f77bea62d479cfd70ec8073692ca911fe306cd5bcf8d6d0000000000000000000000000000000000000000000000000000000000100e58c41df85f0fb47f78779e68b0a0dbefe8d626446b286aebf96741ae274cf3a49d00000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001e0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800fb2e5752573270cb04af9a1ebafc82b67f09a7408217b3b3cb81fe24eb0912900000000000000000000000000000000000000000000000000000000637ce26100000000000000000000000000000000000000000000000000000000002bbecc00000000000000000000000000000000000000000000003cb243ded70e1d0000000000000000000000000000000000000000000000000000000000000000000243d68b3eda5fb3a526daffdab2bf9978f89a6c1a5de19020fd047b80692b67a2712668bc873498a2a69f38ea7874f7e3511baa0637af1b1b304e4cae2bf87c1400000000000000000000000000000000000000000000000000000000000000021e1e506899f5ea70c67ea458042e08a9b45f888f67fa31ae286d760e1c967d14210469b4efc32630d4af00842811b739d7816439abbc9ee48f2d463f3f657fd7",
	//	},
	//}
	mb := make([][]byte, len(m.ChainlinkBlobs))
	for i, blob := range m.ChainlinkBlobs {
		blobBytes, err := hexutil.Decode(blob)
		if err != nil {
			return [][]byte{}, err
		}
		mb[i] = blobBytes

	}
	return mb, nil
}

func generateHMAC(method string, path string, body []byte, clientId string, secret string) string {
	bodyHash := sha256.New()
	bodyHash.Write(body)
	hashString := fmt.Sprintf("%s %s %s %s %d",
		method,
		path,
		hex.EncodeToString(bodyHash.Sum(nil)),
		clientId,
		time.Now().Unix())
	signedMessage := hmac.New(sha256.New, []byte(secret))
	signedMessage.Write([]byte(hashString))
	userHmac := hex.EncodeToString(signedMessage.Sum(nil))
	return userHmac
}

// setCachesOnAPIErr when an off chain look up request fails or gets a 4xx/5xx response code we increment error count and put the upkeep in cooldown state
func (r *EvmRegistry) setCachesOnAPIErr(upkeepId *big.Int) {
	errCount := 1
	cacheKey := upkeepId.String()
	e, ok := r.mercury.apiErrCache.Get(cacheKey)
	if ok {
		errCount = e.(int) + 1
	}

	// With a 10m Error Cache Window every error sets the error count and resets the TTL to 10m
	// On every error that hits during this rolling 10m the error count is increased and the cooldown period by associate is increased
	// This means the user will suffer a max cooldown of 17m on the 10th error at which point the error cache will have expired since its window is 10m
	// After that the user will reset to 0 and start over after a combined total of 34m in cooldown state.

	// increment error count and reset expiration to shift window with last seen error
	r.mercury.apiErrCache.Set(cacheKey, errCount, cache.DefaultExpiration)
	// put upkeep in cooldown state for 2^errors seconds.
	r.mercury.cooldownCache.Set(cacheKey, nil, time.Second*time.Duration(2^errCount))
}
