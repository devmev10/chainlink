import { ethers } from 'hardhat'
import { assert, expect } from 'chai'
import { evmRevert } from '../../test-helpers/matchers'
import { getUsers, Personas } from '../../test-helpers/setup'
import { BigNumber, Signer } from 'ethers'
import { LinkToken__factory as LinkTokenFactory } from '../../../typechain/factories/LinkToken__factory'

import { MockV3Aggregator__factory as MockV3AggregatorFactory } from '../../../typechain/factories/MockV3Aggregator__factory'
import { UpkeepMock__factory as UpkeepMockFactory } from '../../../typechain/factories/UpkeepMock__factory'
import { KeeperRegistrar21 as KeeperRegistrar } from '../../../typechain/KeeperRegistrar21'
import { IKeeperRegistryMaster as IKeeperRegistry } from '../../../typechain/IKeeperRegistryMaster'
import { IKeeperRegistryMaster__factory as IKeeperRegistryMasterFactory } from '../../../typechain/factories/IKeeperRegistryMaster__factory'
import { KeeperRegistry21__factory as KeeperRegistryFactory } from '../../../typechain/factories/KeeperRegistry21__factory'
import { KeeperRegistrar21__factory as KeeperRegistrarFactory } from '../../../typechain/factories/KeeperRegistrar21__factory'

import { MockV3Aggregator } from '../../../typechain/MockV3Aggregator'
import { LinkToken } from '../../../typechain/LinkToken'
import { UpkeepMock } from '../../../typechain/UpkeepMock'
import { toWei } from '../../test-helpers/helpers'

let linkTokenFactory: LinkTokenFactory
let mockV3AggregatorFactory: MockV3AggregatorFactory
let keeperRegistryFactory: KeeperRegistryFactory
let keeperRegistrar: KeeperRegistrarFactory
let upkeepMockFactory: UpkeepMockFactory

let personas: Personas

before(async () => {
  personas = (await getUsers()).personas

  linkTokenFactory = await ethers.getContractFactory('LinkToken')
  mockV3AggregatorFactory = (await ethers.getContractFactory(
    'src/v0.8/tests/MockV3Aggregator.sol:MockV3Aggregator',
  )) as unknown as MockV3AggregatorFactory
  keeperRegistryFactory = (await ethers.getContractFactory(
    'KeeperRegistry2_1',
  )) as unknown as KeeperRegistryFactory
  keeperRegistrar = (await ethers.getContractFactory(
    'KeeperRegistrar2_1',
  )) as unknown as KeeperRegistrarFactory
  upkeepMockFactory = await ethers.getContractFactory('UpkeepMock')
})

const errorMsgs = {
  onlyOwner: 'revert Only callable by owner',
  onlyAdmin: 'OnlyAdminOrOwner()',
  hashPayload: 'HashMismatch()',
  requestNotFound: 'RequestNotFound()',
}

describe('KeeperRegistrar2_1', () => {
  const upkeepName = 'SampleUpkeep'

  const linkEth = BigNumber.from(300000000)
  const gasWei = BigNumber.from(100)
  const executeGas = BigNumber.from(100000)
  const paymentPremiumPPB = BigNumber.from(250000000)
  const flatFeeMicroLink = BigNumber.from(0)
  const maxAllowedAutoApprove = 5
  const offchainConfig = '0x01234567'

  const emptyBytes = '0x00'
  const stalenessSeconds = BigNumber.from(43820)
  const gasCeilingMultiplier = BigNumber.from(1)
  const checkGasLimit = BigNumber.from(20000000)
  const fallbackGasPrice = BigNumber.from(200)
  const fallbackLinkPrice = BigNumber.from(200000000)
  const maxCheckDataSize = BigNumber.from(10000)
  const maxPerformDataSize = BigNumber.from(10000)
  const maxPerformGas = BigNumber.from(5000000)
  const minUpkeepSpend = BigNumber.from('1000000000000000000')
  const amount = BigNumber.from('5000000000000000000')
  const amount1 = BigNumber.from('6000000000000000000')
  const transcoder = ethers.constants.AddressZero

  // Enum values are not auto exported in ABI so have to manually declare
  const autoApproveType_DISABLED = 0
  const autoApproveType_ENABLED_SENDER_ALLOWLIST = 1
  const autoApproveType_ENABLED_ALL = 2

  let owner: Signer
  let admin: Signer
  let someAddress: Signer
  let registrarOwner: Signer
  let stranger: Signer
  let requestSender: Signer

  let linkToken: LinkToken
  let linkEthFeed: MockV3Aggregator
  let gasPriceFeed: MockV3Aggregator
  let registry: IKeeperRegistry
  let mock: UpkeepMock
  let registrar: KeeperRegistrar

  beforeEach(async () => {
    owner = personas.Default
    admin = personas.Neil
    someAddress = personas.Ned
    registrarOwner = personas.Nelly
    stranger = personas.Nancy
    requestSender = personas.Norbert

    linkToken = await linkTokenFactory.connect(owner).deploy()
    gasPriceFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(0, gasWei)
    linkEthFeed = await mockV3AggregatorFactory
      .connect(owner)
      .deploy(9, linkEth)

    registry = IKeeperRegistryMasterFactory.connect(
      (
        await keeperRegistryFactory
          .connect(owner)
          .deploy(
            0,
            linkToken.address,
            linkEthFeed.address,
            gasPriceFeed.address,
          )
      ).address,
      owner,
    )

    mock = await upkeepMockFactory.deploy()

    registrar = await keeperRegistrar
      .connect(registrarOwner)
      .deploy(
        linkToken.address,
        autoApproveType_DISABLED,
        BigNumber.from('0'),
        registry.address,
        minUpkeepSpend,
      )

    await linkToken
      .connect(owner)
      .transfer(await requestSender.getAddress(), toWei('1000'))

    const keepers = [
      await personas.Carol.getAddress(),
      await personas.Nancy.getAddress(),
      await personas.Ned.getAddress(),
      await personas.Neil.getAddress(),
    ]
    const config = {
      paymentPremiumPPB,
      flatFeeMicroLink,
      checkGasLimit,
      stalenessSeconds,
      gasCeilingMultiplier,
      minUpkeepSpend,
      maxCheckDataSize,
      maxPerformDataSize,
      maxPerformGas,
      fallbackGasPrice,
      fallbackLinkPrice,
      transcoder,
      registrar: registrar.address,
    }
    const onchainConfig = ethers.utils.defaultAbiCoder.encode(
      [
        'tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds\
          ,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,\
          uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,\
          address registrar)',
      ],
      [config],
    )
    await registry
      .connect(owner)
      .setConfig(keepers, keepers, 1, onchainConfig, 1, '0x')
  })

  describe('#typeAndVersion', () => {
    it('uses the correct type and version', async () => {
      const typeAndVersion = await registrar.typeAndVersion()
      assert.equal(typeAndVersion, 'KeeperRegistrar 2.1.0')
    })
  })

  describe('#register', () => {
    it('reverts if not called by the LINK token', async () => {
      await evmRevert(
        registrar
          .connect(someAddress)
          .register(
            upkeepName,
            emptyBytes,
            mock.address,
            executeGas,
            await admin.getAddress(),
            emptyBytes,
            offchainConfig,
            amount,
            await requestSender.getAddress(),
          ),
        'OnlyLink()',
      )
    })

    it('reverts if the amount passed in data mismatches actual amount sent', async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount1,
          await requestSender.getAddress(),
        ],
      )

      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'AmountMismatch()',
      )
    })

    it('reverts if the sender passed in data mismatches actual sender', async () => {
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await admin.getAddress(), // Should have been requestSender.getAddress()
        ],
      )
      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'SenderMismatch()',
      )
    })

    it('reverts if the admin address is 0x0000...', async () => {
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          '0x0000000000000000000000000000000000000000',
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )

      await evmRevert(
        linkToken
          .connect(requestSender)
          .transferAndCall(registrar.address, amount, abiEncodedBytes),
        'RegistrationRequestFailed()',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve ON
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)

      const [id] = await registry.getActiveUpkeepIDs(0, 1)

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.executeGas, executeGas.toNumber())
      assert.equal(newupkeep.offchainConfig, offchainConfig)

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve OFF - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      //get upkeep count before attempting registration
      const beforeCount = (await registry.getState()).state.numUpkeeps

      //set auto approve OFF, threshold limits dont matter in this case
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()

      //get upkeep count after attempting registration
      const afterCount = (await registry.getState()).state.numUpkeeps
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount)

      //confirm that only RegistrationRequested event is emitted and RegistrationApproved event is not
      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).not.to.emit(registrar, 'RegistrationApproved')

      const hash = receipt.logs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(await admin.getAddress(), pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    })

    it('Auto Approve ON - Throttle max approvals - does not register an upkeep on KeeperRegistry beyond the max limit, emits only RegistrationRequested event after limit is hit', async () => {
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 0)

      //set auto approve on, with max 1 allowed
      await registrar.connect(registrarOwner).setRegistrationConfig(
        autoApproveType_ENABLED_ALL,
        1, // maxAllowedAutoApprove
        registry.address,
        minUpkeepSpend,
      )

      //register within threshold, new upkeep should be registered
      let abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas,
        await admin.getAddress(),
        emptyBytes,
        offchainConfig,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      //try registering another one, new upkeep should not be registered
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 1, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        offchainConfig,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // Still 1

      // Now set new max limit to 2. One more upkeep should get auto approved
      await registrar.connect(registrarOwner).setRegistrationConfig(
        autoApproveType_ENABLED_ALL,
        2, // maxAllowedAutoApprove
        registry.address,
        minUpkeepSpend,
      )
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 2, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        offchainConfig,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 2) // 1 -> 2

      // One more upkeep should not get registered
      abiEncodedBytes = registrar.interface.encodeFunctionData('register', [
        upkeepName,
        emptyBytes,
        mock.address,
        executeGas.toNumber() + 3, // make unique hash
        await admin.getAddress(),
        emptyBytes,
        offchainConfig,
        amount,
        await requestSender.getAddress(),
      ])
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 2) // Still 2
    })

    it('Auto Approve Sender Allowlist - sender in allowlist - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // Add sender to allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, true)

      //register with auto approve ON
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)

      const [id] = await registry.getActiveUpkeepIDs(0, 1)

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.executeGas, executeGas.toNumber())

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('Auto Approve Sender Allowlist - sender NOT in allowlist - does not registers an upkeep on KeeperRegistry, emits only RegistrationRequested event', async () => {
      const beforeCount = (await registry.getState()).state.numUpkeeps
      const senderAddress = await requestSender.getAddress()

      //set auto approve to ENABLED_SENDER_ALLOWLIST type with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_SENDER_ALLOWLIST,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // Explicitly remove sender from allowlist
      await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, false)

      //register. auto approve shouldn't happen
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()

      //get upkeep count after attempting registration
      const afterCount = (await registry.getState()).state.numUpkeeps
      //confirm that a new upkeep has NOT been registered and upkeep count is still the same
      assert.deepEqual(beforeCount, afterCount)

      //confirm that only RegistrationRequested event is emitted and RegistrationApproved event is not
      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).not.to.emit(registrar, 'RegistrationApproved')

      const hash = receipt.logs[2].topics[1]
      const pendingRequest = await registrar.getPendingRequest(hash)
      assert.equal(await admin.getAddress(), pendingRequest[0])
      assert.ok(amount.eq(pendingRequest[1]))
    })
  })

  describe('#registerUpkeep', () => {
    it('reverts with empty message if amount sent is not available in LINK allowance', async () => {
      await evmRevert(
        registrar.connect(someAddress).registerUpkeep({
          name: upkeepName,
          upkeepContract: mock.address,
          gasLimit: executeGas,
          adminAddress: await admin.getAddress(),
          checkData: emptyBytes,
          offchainConfig: emptyBytes,
          amount,
          encryptedEmail: emptyBytes,
        }),
        '',
      )
    })

    it('reverts if the amount passed in data is less than configured minimum', async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      // amt is one order of magnitude less than minUpkeepSpend
      const amt = BigNumber.from('100000000000000000')

      await evmRevert(
        registrar.connect(someAddress).registerUpkeep({
          name: upkeepName,
          upkeepContract: mock.address,
          gasLimit: executeGas,
          adminAddress: await admin.getAddress(),
          checkData: emptyBytes,
          offchainConfig: emptyBytes,
          amount: amt,
          encryptedEmail: emptyBytes,
        }),
        'InsufficientPayment()',
      )
    })

    it('Auto Approve ON - registers an upkeep on KeeperRegistry instantly and emits both RegistrationRequested and RegistrationApproved events', async () => {
      //set auto approve ON with high threshold limits
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_ENABLED_ALL,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      await linkToken.connect(requestSender).approve(registrar.address, amount)

      const params = {
        name: upkeepName,
        upkeepContract: mock.address,
        gasLimit: executeGas,
        adminAddress: await admin.getAddress(),
        checkData: emptyBytes,
        offchainConfig,
        amount,
        encryptedEmail: emptyBytes,
      }

      // simulate tx to check return values
      const [upkeepID, forwarder] = await registrar
        .connect(requestSender)
        .callStatic.registerUpkeep(params)
      expect(upkeepID).to.not.equal(0)
      expect(forwarder).to.not.equal(ethers.constants.AddressZero)

      const tx = await registrar.connect(requestSender).registerUpkeep(params)
      assert.equal((await registry.getState()).state.numUpkeeps.toNumber(), 1) // 0 -> 1

      //confirm if a new upkeep has been registered and the details are the same as the one just registered
      const [id] = await registry.getActiveUpkeepIDs(0, 1)
      const newupkeep = await registry.getUpkeep(id)
      assert.equal(newupkeep.target, mock.address)
      assert.equal(newupkeep.admin, await admin.getAddress())
      assert.equal(newupkeep.checkData, emptyBytes)
      assert.equal(newupkeep.balance.toString(), amount.toString())
      assert.equal(newupkeep.executeGas, executeGas.toNumber())
      assert.equal(newupkeep.offchainConfig, offchainConfig)

      await expect(tx).to.emit(registrar, 'RegistrationRequested')
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })
  })

  describe('#setAutoApproveAllowedSender', () => {
    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .setAutoApproveAllowedSender(await admin.getAddress(), false)
      await evmRevert(tx, 'Only callable by owner')
    })

    it('sets the allowed status correctly and emits log', async () => {
      const senderAddress = await stranger.getAddress()
      let tx = await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, true)
      await expect(tx)
        .to.emit(registrar, 'AutoApproveAllowedSenderSet')
        .withArgs(senderAddress, true)

      let senderAllowedStatus = await registrar
        .connect(owner)
        .getAutoApproveAllowedSender(senderAddress)
      assert.isTrue(senderAllowedStatus)

      tx = await registrar
        .connect(registrarOwner)
        .setAutoApproveAllowedSender(senderAddress, false)
      await expect(tx)
        .to.emit(registrar, 'AutoApproveAllowedSenderSet')
        .withArgs(senderAddress, false)

      senderAllowedStatus = await registrar
        .connect(owner)
        .getAutoApproveAllowedSender(senderAddress)
      assert.isFalse(senderAllowedStatus)
    })
  })

  describe('#approve', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )

      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
    })

    it('reverts if not called by the owner', async () => {
      const tx = registrar
        .connect(stranger)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, 'Only callable by owner')
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          emptyBytes,
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })

    it('reverts if any member of the payload changes', async () => {
      let tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          ethers.Wallet.createRandom().address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          10000,
          await admin.getAddress(),
          emptyBytes,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          ethers.Wallet.createRandom().address,
          emptyBytes,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          '0x1234',
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.hashPayload)
    })

    it('approves an existing registration request', async () => {
      const tx = await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          hash,
        )
      await expect(tx).to.emit(registrar, 'RegistrationApproved')
    })

    it('deletes the request afterwards / reverts if the request DNE', async () => {
      await registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          hash,
        )
      const tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })

  describe('#cancel', () => {
    let hash: string

    beforeEach(async () => {
      await registrar
        .connect(registrarOwner)
        .setRegistrationConfig(
          autoApproveType_DISABLED,
          maxAllowedAutoApprove,
          registry.address,
          minUpkeepSpend,
        )

      //register with auto approve OFF
      const abiEncodedBytes = registrar.interface.encodeFunctionData(
        'register',
        [
          upkeepName,
          emptyBytes,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          offchainConfig,
          amount,
          await requestSender.getAddress(),
        ],
      )
      const tx = await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
      const receipt = await tx.wait()
      hash = receipt.logs[2].topics[1]
      // submit duplicate request (increase balance)
      await linkToken
        .connect(requestSender)
        .transferAndCall(registrar.address, amount, abiEncodedBytes)
    })

    it('reverts if not called by the admin / owner', async () => {
      const tx = registrar.connect(stranger).cancel(hash)
      await evmRevert(tx, errorMsgs.onlyAdmin)
    })

    it('reverts if the hash does not exist', async () => {
      const tx = registrar
        .connect(registrarOwner)
        .cancel(
          '0x000000000000000000000000322813fd9a801c5507c9de605d63cea4f2ce6c44',
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })

    it('refunds the total request balance to the admin address if owner cancels', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(registrarOwner).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(2))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('refunds the total request balance to the admin address if admin cancels', async () => {
      const before = await linkToken.balanceOf(await admin.getAddress())
      const tx = await registrar.connect(admin).cancel(hash)
      const after = await linkToken.balanceOf(await admin.getAddress())
      assert.isTrue(after.sub(before).eq(amount.mul(BigNumber.from(2))))
      await expect(tx).to.emit(registrar, 'RegistrationRejected')
    })

    it('deletes the request hash', async () => {
      await registrar.connect(registrarOwner).cancel(hash)
      let tx = registrar.connect(registrarOwner).cancel(hash)
      await evmRevert(tx, errorMsgs.requestNotFound)
      tx = registrar
        .connect(registrarOwner)
        .approve(
          upkeepName,
          mock.address,
          executeGas,
          await admin.getAddress(),
          emptyBytes,
          emptyBytes,
          hash,
        )
      await evmRevert(tx, errorMsgs.requestNotFound)
    })
  })
})
