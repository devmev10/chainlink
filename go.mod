module github.com/smartcontractkit/chainlink/v2

go 1.20

require (
	github.com/CosmWasm/wasmd v0.30.0
	github.com/Depado/ginprom v1.7.11
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/ava-labs/coreth v0.12.1
	github.com/btcsuite/btcd v0.23.4
	github.com/cosmos/cosmos-sdk v0.47.2
	github.com/danielkov/gin-helmet v0.0.0-20171108135313-1387e224435e
	github.com/ethereum/go-ethereum v1.11.6
	github.com/fatih/color v1.15.0
	github.com/fxamacker/cbor/v2 v2.4.0
	github.com/gagliardetto/solana-go v1.8.2
	github.com/getsentry/sentry-go v0.21.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-contrib/expvar v0.0.1
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-contrib/size v0.0.0-20230212012657-e14a14094dc4
	github.com/gin-gonic/gin v1.9.0
	github.com/go-webauthn/webauthn v0.8.2
	github.com/gogo/protobuf v1.3.3
	github.com/google/pprof v0.0.0-20230510103437-eeec1cb781c3
	github.com/google/uuid v1.3.0
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.5.0
	github.com/graph-gophers/dataloader v5.0.0+incompatible
	github.com/graph-gophers/graphql-go v1.5.0
	github.com/hashicorp/go-plugin v1.4.9
	github.com/hdevalence/ed25519consensus v0.1.0
	github.com/jackc/pgconn v1.14.0
	github.com/jackc/pgtype v1.14.0
	github.com/jackc/pgx/v4 v4.18.1
	github.com/jpillora/backoff v1.0.0
	github.com/kylelemons/godebug v1.1.0
	github.com/leanovate/gopter v0.2.10-0.20210127095200-9abe2343507a
	github.com/lib/pq v1.10.9
	github.com/libp2p/go-libp2p-core v0.20.1
	github.com/libp2p/go-libp2p-peerstore v0.8.0
	github.com/manyminds/api2go v0.0.0-20220325145637-95b4fb838cf6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multiaddr v0.9.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/gomega v1.27.6
	github.com/pelletier/go-toml v1.9.5
	github.com/pelletier/go-toml/v2 v2.0.7
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.11.2
	github.com/prometheus/client_golang v1.15.1
	github.com/prometheus/client_model v0.4.0
	github.com/pyroscope-io/client v0.7.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rogpeppe/go-internal v1.10.0
	github.com/scylladb/go-reflectx v1.0.1
	github.com/shirou/gopsutil/v3 v3.23.4
	github.com/shopspring/decimal v1.3.1
	github.com/smartcontractkit/chainlink-cosmos v0.4.0
	github.com/smartcontractkit/chainlink-relay v0.1.7-0.20230505214134-c890447508f9
	github.com/smartcontractkit/chainlink-solana v1.0.3-0.20230424191709-c9fec8c08e1b
	github.com/smartcontractkit/chainlink-starknet/relayer v0.0.0-20230510171651-62cd01eeb9f8
	github.com/smartcontractkit/libocr v0.0.0-20230503222226-29f534b2de1a
	github.com/smartcontractkit/ocr2keepers v0.6.15
	github.com/smartcontractkit/ocr2vrf v0.0.0-20230510102715-c58be582bf19
	github.com/smartcontractkit/sqlx v1.3.5-0.20210805004948-4be295aacbeb
	github.com/smartcontractkit/wsrpc v0.7.1
	github.com/spf13/cast v1.5.0
	github.com/stretchr/testify v1.8.2
	github.com/tendermint/tendermint v0.35.9
	github.com/theodesp/go-heaps v0.0.0-20190520121037-88e35354fe0a
	github.com/tidwall/gjson v1.14.4
	github.com/ugorji/go/codec v1.2.11
	github.com/ulule/limiter/v3 v3.11.1
	github.com/umbracle/ethgo v0.1.3
	github.com/unrolled/secure v1.13.0
	github.com/urfave/cli v1.22.13
	go.dedis.ch/fixbuf v1.0.3
	go.dedis.ch/kyber/v3 v3.1.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.9.0
	golang.org/x/exp v0.0.0-20230510235704-dd950f8aeaea
	golang.org/x/sync v0.2.0
	golang.org/x/term v0.8.0
	golang.org/x/text v0.9.0
	golang.org/x/tools v0.9.1
	gonum.org/v1/gonum v0.13.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/guregu/null.v2 v2.1.2
	gopkg.in/guregu/null.v4 v4.0.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14 // indirect
	cosmossdk.io/api v0.4.1 // indirect
	cosmossdk.io/core v0.6.1 // indirect
	cosmossdk.io/depinject v1.0.0-alpha.3 // indirect
	cosmossdk.io/errors v1.0.0-beta.7 // indirect
	cosmossdk.io/math v1.0.0 // indirect
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/ChainSafe/go-schnorrkel v1.0.0 // indirect
	github.com/CosmWasm/wasmvm v1.2.3 // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/NethermindEth/juno v0.0.0-20220630151419-cbd368b222ac // indirect
	github.com/VictoriaMetrics/fastcache v1.12.1 // indirect
	github.com/andres-erbsen/clock v0.0.0-20160526145045-9e14626cd129 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/aybabtme/rgbterm v0.0.0-20170906152045-cc83f3b3ce59 // indirect
	github.com/benbjohnson/clock v1.3.4 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.1-0.20220910012023-760eaf8b6816 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2 // indirect
	github.com/bytedance/sonic v1.8.8 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230510135629-fe7ae7a62e0f // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/cometbft/cometbft v0.37.1 // indirect
	github.com/cometbft/cometbft-db v0.8.0 // indirect
	github.com/confio/ics23/go v0.9.0 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-db v1.0.0-rc.1 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.3 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogoproto v1.4.9 // indirect
	github.com/cosmos/gorocksdb v1.2.0 // indirect
	github.com/cosmos/iavl v0.21.0 // indirect
	github.com/cosmos/ibc-go/v4 v4.3.0 // indirect
	github.com/cosmos/ics23/go v0.10.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.0 // indirect
	github.com/cosmos/ledger-go v0.9.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/deckarep/golang-set/v2 v2.3.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dfuse-io/logging v0.0.0-20210109005628-b97a57253f70 // indirect
	github.com/dgraph-io/badger/v2 v2.2007.4 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dontpanicdao/caigo v0.4.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/elastic/gosigar v0.14.2 // indirect
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/flynn/noise v1.0.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gagliardetto/binary v0.7.8 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/gedex/inflector v0.0.0-20170307190818-16278e9db813 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-kit/kit v0.12.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.13.0 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/go-webauthn/revoke v0.1.9 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/glog v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/go-tpm v0.3.3 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-bexpr v0.1.10 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.2 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/huin/goupnp v1.2.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ipfs/boxo v0.8.1 // indirect
	github.com/ipfs/go-cid v0.4.1 // indirect
	github.com/ipfs/go-datastore v0.6.0 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipns v0.3.0 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-log/v2 v2.5.1 // indirect
	github.com/ipld/go-ipld-prime v0.20.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/go-temp-err-catcher v0.1.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/koron/go-ssdp v0.0.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/libp2p/go-addr-util v0.2.0 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-conn-security-multistream v0.4.0 // indirect
	github.com/libp2p/go-eventbus v0.3.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p v0.27.3 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.3.0 // indirect
	github.com/libp2p/go-libp2p-autonat v0.8.0 // indirect
	github.com/libp2p/go-libp2p-blankhost v0.4.0 // indirect
	github.com/libp2p/go-libp2p-circuit v0.6.0 // indirect
	github.com/libp2p/go-libp2p-discovery v0.7.0 // indirect
	github.com/libp2p/go-libp2p-kad-dht v0.23.0 // indirect
	github.com/libp2p/go-libp2p-kbucket v0.5.0 // indirect
	github.com/libp2p/go-libp2p-loggables v0.2.0 // indirect
	github.com/libp2p/go-libp2p-mplex v0.8.0 // indirect
	github.com/libp2p/go-libp2p-nat v0.2.0 // indirect
	github.com/libp2p/go-libp2p-noise v0.5.0 // indirect
	github.com/libp2p/go-libp2p-pnet v0.3.0 // indirect
	github.com/libp2p/go-libp2p-record v0.2.0 // indirect
	github.com/libp2p/go-libp2p-swarm v0.11.0 // indirect
	github.com/libp2p/go-libp2p-tls v0.5.0 // indirect
	github.com/libp2p/go-libp2p-transport-upgrader v0.8.0 // indirect
	github.com/libp2p/go-libp2p-yamux v0.10.0 // indirect
	github.com/libp2p/go-mplex v0.7.0 // indirect
	github.com/libp2p/go-msgio v0.3.0 // indirect
	github.com/libp2p/go-nat v0.1.0 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
	github.com/libp2p/go-openssl v0.1.0 // indirect
	github.com/libp2p/go-reuseport v0.3.0 // indirect
	github.com/libp2p/go-reuseport-transport v0.2.0 // indirect
	github.com/libp2p/go-sockaddr v0.2.0 // indirect
	github.com/libp2p/go-stream-muxer-multistream v0.5.0 // indirect
	github.com/libp2p/go-tcp-transport v0.6.1 // indirect
	github.com/libp2p/go-ws-transport v0.7.0 // indirect
	github.com/libp2p/go-yamux/v2 v2.3.0 // indirect
	github.com/libp2p/go-yamux/v4 v4.0.0 // indirect
	github.com/linxGnu/grocksdb v1.7.16 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.54 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/mimoo/StrobeGo v0.0.0-20220103164710-9a04d6ca976b // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/pointerstructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr-dns v0.3.1 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multiaddr-net v0.2.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.9.0 // indirect
	github.com/multiformats/go-multihash v0.2.1 // indirect
	github.com/multiformats/go-multistream v0.4.1 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230110094441-db37f07504ce // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/onsi/ginkgo/v2 v2.9.4 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/opencontainers/runc v1.1.7 // indirect
	github.com/opencontainers/runtime-spec v1.1.0-rc.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/petermattis/goid v0.0.0-20230317030725-371a4b8eda08 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/prometheus/common v0.43.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/pyroscope-io/godeltaprof v0.1.1 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.3.2 // indirect
	github.com/quic-go/qtls-go1-20 v0.2.2 // indirect
	github.com/quic-go/quic-go v0.34.0 // indirect
	github.com/quic-go/webtransport-go v0.5.2 // indirect
	github.com/raulk/go-watchdog v1.3.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/regen-network/cosmos-proto v0.3.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/rs/zerolog v1.29.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/smartcontractkit/caigo v0.0.0-20230508053235-41120ca1f9f3 // indirect
	github.com/spacemonkeygo/spacelog v0.0.0-20180420211403-2296661a0572 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cobra v1.7.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	github.com/status-im/keycard-go v0.2.0 // indirect
	github.com/streamingfast/logging v0.0.0-20221209193439-bff11742bf4c // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tecbot/gorocksdb v0.0.0-20191217155057-f0fad39f321c // indirect
	github.com/tendermint/btcd v0.1.1 // indirect
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15 // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tendermint/tm-db v0.6.7 // indirect
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569 // indirect
	github.com/tidwall/btree v1.6.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/umbracle/fastrlp v0.0.0-20220527094140-59d5dd30e722 // indirect
	github.com/urfave/cli/v2 v2.17.2-0.20221006022127-8f469abc00aa // indirect
	github.com/valyala/fastjson v1.4.1 // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zondax/hid v0.9.1 // indirect
	github.com/zondax/ledger-go v0.14.1 // indirect
	go.dedis.ch/protobuf v1.0.11 // indirect
	go.etcd.io/bbolt v1.3.7 // indirect
	go.mongodb.org/mongo-driver v1.11.6 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.15.1 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/fx v1.19.3 // indirect
	go.uber.org/ratelimit v0.2.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
	pgregory.net/rapid v0.5.7 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	// Fix go mod tidy issue for ambiguous imports from go-ethereum
	// See https://github.com/ugorji/go/issues/279
	github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.1

	// replicating the replace directive on cosmos SDK
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
