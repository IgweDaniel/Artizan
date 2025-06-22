import { expect } from "chai";
import hre from "hardhat";
import { AbiCoder, TypedDataDomain, ZeroAddress } from "ethers";
import { randomBytes } from "crypto";

describe("LazyMintZone", function () {
  async function deployFixture() {
    const [owner, signer, account1] = await hre.ethers.getSigners();
    const LazyMint1155 = await hre.ethers.getContractFactory("LazyMint1155");
    const lazyMint1155 = await LazyMint1155.deploy(signer.address);
    const LazyMintZone = await hre.ethers.getContractFactory("LazyMintZone");
    const lazyMintZone = await LazyMintZone.deploy(
      owner.address,
      await lazyMint1155.getAddress()
    );
    return { lazyMintZone, lazyMint1155, owner, signer, account1 };
  }

  describe("Deployment", function () {
    it("Should set the correct NFT contract address", async function () {
      const { lazyMintZone, lazyMint1155 } = await deployFixture();
      expect(await lazyMintZone.nft()).to.equal(
        await lazyMint1155.getAddress()
      );
    });
    it("Should set the correct owner", async function () {
      const { lazyMintZone, owner } = await deployFixture();
      expect(await lazyMintZone.owner()).to.equal(owner.address);
    });
  });

  describe("setNftAddress", function () {
    it("Should allow the owner to set a new NFT address", async function () {
      const { lazyMintZone, owner } = await deployFixture();
      const newNft = hre.ethers.Wallet.createRandom().address;
      await lazyMintZone.connect(owner).setNftAddress(newNft);
      expect(await lazyMintZone.nft()).to.equal(newNft);
    });
    it("Should revert if a non-owner tries to set the NFT address", async function () {
      const { lazyMintZone, signer } = await deployFixture();
      const newNft = hre.ethers.Wallet.createRandom().address;
      await expect(
        lazyMintZone.connect(signer).setNftAddress(newNft)
      ).to.be.revertedWithCustomError(
        lazyMintZone,
        "OwnableUnauthorizedAccount"
      );
    });
    it("Should revert if the new NFT address is zero", async function () {
      const { lazyMintZone, owner } = await deployFixture();
      await expect(
        lazyMintZone
          .connect(owner)
          .setNftAddress("0x0000000000000000000000000000000000000000")
      ).to.be.revertedWith("Invalid NFT address");
    });
  });

  describe("getSeaportMetadata", function () {
    it("Should return the correct name and schema", async function () {
      const { lazyMintZone } = await deployFixture();
      const [name, schemas] = await lazyMintZone.getSeaportMetadata();
      expect(name).to.equal("ArtiartZone");
      expect(schemas[0].id).to.equal(3003);
    });
  });

  describe("supportsInterface", function () {
    it("Should return true for ZoneInterface and ERC165", async function () {
      const { lazyMintZone } = await deployFixture();

      expect(await lazyMintZone.supportsInterface("0x39dd6933")).to.be.true;
      expect(await lazyMintZone.supportsInterface("0x01ffc9a7")).to.be.true;
    });
    it("Should return false for random interface IDs", async function () {
      const { lazyMintZone } = await deployFixture();
      expect(await lazyMintZone.supportsInterface("0xffffffff")).to.be.false;
    });
  });

  describe("authorizeOrder", function () {
    it("Should mint the NFT if not already minted, using a valid voucher in extraData", async function () {
      const { lazyMintZone, lazyMint1155, signer } = await deployFixture();
      // Prepare a valid voucher
      const tokenId = 123;
      const amount = 1;
      const uri = "ipfs://test";
      const domain: TypedDataDomain = {
        name: "LazyMint1155",
        version: "1",
        chainId: hre.network.config.chainId,
        verifyingContract: await lazyMint1155.getAddress(),
      };
      const types = {
        Voucher: [
          { name: "owner", type: "address" },
          { name: "tokenId", type: "uint256" },
          { name: "amount", type: "uint256" },
          { name: "uri", type: "string" },
        ],
      };
      const voucher = { owner: signer.address, tokenId, amount, uri };
      const signature = await signer.signTypedData(domain, types, voucher);
      const abi = AbiCoder.defaultAbiCoder();
      const extraData = abi.encode(
        [
          "tuple(address owner,uint256 tokenId,uint256 amount,string uri,bytes signature)",
        ],
        [[signer.address, tokenId, amount, uri, signature]]
      );

      // Mock all required ZoneParameters fields, including orderHashes
      const zoneParameters = {
        orderHash: "0x" + "00".repeat(32),
        orderHashes: [], // Added missing field
        fulfiller: signer.address,
        offer: [],
        consideration: [],
        extraData,
        orderType: 0,
        zoneHash: "0x" + "00".repeat(32),
        startTime: 0,
        endTime: 0,
        zone: lazyMintZone.target,
        offerer: signer.address,
        conduitKey: "0x" + "00".repeat(32),
        totalOriginalConsiderationItems: 0,
      };
      // Call authorizeOrder
      await lazyMintZone.authorizeOrder(zoneParameters);
      const returnValue = await lazyMintZone.authorizeOrder.staticCall(
        zoneParameters
      );

      const expectedSelector =
        lazyMintZone.interface.getFunction("authorizeOrder").selector;

      expect(returnValue).to.equal(expectedSelector);
      expect(await lazyMint1155.isTokenMinted(tokenId)).to.be.true;
    });

    it("Should revert if the voucher is invalid", async function () {
      const { lazyMintZone, signer, account1 } = await deployFixture();
      // Prepare an invalid voucher
      const tokenId = 123;
      const amount = 1;
      const uri = "ipfs://test";

      const abi = AbiCoder.defaultAbiCoder();
      const extraData = abi.encode(
        ["tuple(address owner,uint256 tokenId,uint256 amount,string uri)"],
        [[signer.address, tokenId, amount, uri]]
      );

      // Mock all required ZoneParameters fields, including orderHashes
      const zoneParameters = {
        orderHash: "0x" + "00".repeat(32),
        orderHashes: [], // Added missing field
        fulfiller: signer.address,
        offer: [],
        consideration: [],
        extraData,
        orderType: 0,
        zoneHash: "0x" + "00".repeat(32),
        startTime: 0,
        endTime: 0,
        zone: lazyMintZone.target,
        offerer: signer.address,
        conduitKey: "0x" + "00".repeat(32),
        totalOriginalConsiderationItems: 0,
      };
      // Call authorizeOrder and expect it to revert
      await expect(
        lazyMintZone.authorizeOrder(zoneParameters)
      ).to.be.revertedWithPanic();
    });
    // Additional negative tests can be added here for invalid voucher, already minted, etc.
  });

  describe("validateOrder", function () {
    it("Should return the correct selector", async function () {
      const { lazyMintZone } = await deployFixture();

      const zoneParameters = {
        orderHash: "0x" + "00".repeat(32),
        orderHashes: [], // Added missing field
        fulfiller: ZeroAddress,
        offer: [],
        consideration: [],
        extraData: randomBytes(32),
        orderType: 0,
        zoneHash: "0x" + "00".repeat(32),
        startTime: 0,
        endTime: 0,
        zone: lazyMintZone.target,
        offerer: ZeroAddress,
        conduitKey: "0x" + "00".repeat(32),
        totalOriginalConsiderationItems: 0,
      };
      // Call authorizeOrder
      const returnValue = await lazyMintZone.validateOrder.staticCall(
        zoneParameters
      );

      const expectedSelector =
        lazyMintZone.interface.getFunction("validateOrder").selector;

      expect(returnValue).to.equal(expectedSelector);
    });
  });
});
