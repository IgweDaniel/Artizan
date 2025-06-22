import {
  time,
  loadFixture,
} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import { anyValue } from "@nomicfoundation/hardhat-chai-matchers/withArgs";
import { expect } from "chai";
import hre from "hardhat";
import {
  AbiCoder,
  AddressLike,
  getBytes,
  keccak256,
  TypedDataDomain,
} from "ethers";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

describe("LazyMint1155", function () {
  // We define a fixture to reuse the same setup in every test.
  // We use loadFixture to run this setup once, snapshot that state,
  // and reset Hardhat Network to that snapshot in every test.
  async function deployLazyMint1155Fixture() {
    const [owner, signer, account1] = await hre.ethers.getSigners();

    const LazyMint1155 = await hre.ethers.getContractFactory("LazyMint1155");
    const lazyMint1155 = await LazyMint1155.deploy(signer.address);

    const domain: TypedDataDomain = {
      name: "LazyMint1155",
      version: "1",
      chainId: hre.network.config.chainId,
      verifyingContract: await lazyMint1155.getAddress(),
    };

    return { lazyMint1155, owner, signer, account1, domain };
  }

  describe("Deployment", function () {
    it("Should set the right signer", async function () {
      const { lazyMint1155, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      expect(await (lazyMint1155 as any).signer()).to.equal(signer.address);
    });
    it("Should set the right owner", async function () {
      const { lazyMint1155, owner } = await loadFixture(
        deployLazyMint1155Fixture
      );
      expect(await lazyMint1155.owner()).to.equal(owner.address);
    });
  });

  describe("Set Signer", function () {
    it("Should allow the owner to set a new signer", async function () {
      const { lazyMint1155, owner, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const newSigner = hre.ethers.Wallet.createRandom().address;
      await lazyMint1155.connect(owner).setSigner(newSigner);

      expect(await (lazyMint1155 as any).signer()).to.equal(newSigner);
    });
    it("Should not allow non-owner to set a new signer", async function () {
      const { lazyMint1155, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const newSigner = hre.ethers.Wallet.createRandom().address;
      await expect(
        lazyMint1155.connect(signer).setSigner(newSigner)
      ).to.be.revertedWithCustomError(
        lazyMint1155,
        "OwnableUnauthorizedAccount"
      );
    });
  });

  describe("Mint with Voucher", function () {
    async function createVoucher(
      {
        owner,
        tokenId,
        amount,
        uri,
      }: {
        owner: string;
        tokenId: number;
        amount: number;
        uri: string;
      },
      domain: TypedDataDomain,
      signer: HardhatEthersSigner
    ) {
      const types = {
        Voucher: [
          { name: "owner", type: "address" },
          { name: "tokenId", type: "uint256" },
          { name: "amount", type: "uint256" },
          { name: "uri", type: "string" },
        ],
      };

      const voucher = { owner, tokenId, amount, uri };

      const signature = await signer.signTypedData(domain, types, voucher);

      return {
        owner: owner,
        tokenId: tokenId,
        amount: amount,
        uri: uri,
        signature: signature,
      };
    }

    const tokenId = 189; // Unique ID for this token
    const amount = Math.floor(Math.random() * 50) + 1; // Number of tokens to mint
    const uri = "ipfs://metadata-url";

    it("Should allow lazy minting of tokens and update minted status", async function () {
      const { lazyMint1155, owner, signer, domain } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const recipientAddress = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      // Create the voucher object
      const voucher = await createVoucher(
        {
          owner: recipientAddress,
          tokenId: tokenId,
          amount: amount,
          uri: uri,
        },
        domain,
        signer
      );

      let balance = await lazyMint1155.balanceOf(recipientAddress, tokenId);
      expect(balance).to.equal(0);

      let isMinted = await lazyMint1155.isTokenMinted(tokenId);
      expect(isMinted).to.be.false;

      await lazyMint1155
        .connect(owner)
        .mintIfNotExists(voucher, recipientAddress);

      balance = await lazyMint1155.balanceOf(recipientAddress, tokenId);
      expect(balance).to.equal(amount);

      isMinted = await lazyMint1155.isTokenMinted(tokenId);
      expect(isMinted).to.be.true;
    });

    // test reminting should do nothing if already minted
    it("Should not remint or panic if reminting of already minted tokens", async function () {
      const { lazyMint1155, owner, signer, domain } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const recipientAddress = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      // Create the voucher object
      const voucher = await createVoucher(
        {
          owner: recipientAddress,
          tokenId: tokenId,
          amount: amount,
          uri: uri,
        },
        domain,
        signer
      );

      // The recipient address

      await lazyMint1155
        .connect(owner)
        .mintIfNotExists(voucher, recipientAddress);

      // Attempt to remint the same token
      await lazyMint1155
        .connect(owner)
        .mintIfNotExists(voucher, recipientAddress);

      // Verify balance remains the same (no additional minting)
      const balance = await lazyMint1155.balanceOf(recipientAddress, tokenId);
      expect(balance).to.equal(amount);
    });

    it("Should panic if recipient not voucher address", async function () {
      const { lazyMint1155, owner, domain } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const recipientAddress = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      // Create the voucher object
      const voucher = await createVoucher(
        {
          owner: recipientAddress,
          tokenId: tokenId,
          amount: amount,
          uri: uri,
        },
        domain,
        owner
      );

      // The recipient address

      await expect(
        lazyMint1155
          .connect(owner)
          .mintIfNotExists(
            voucher,
            "0x1234567890123456789012345678901234567890"
          )
      ).to.be.revertedWith("Voucher owner mismatch");
    });
    it("Should panic with invalid signer", async function () {
      const { lazyMint1155, owner, domain } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const recipientAddress = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      // Create the voucher object
      const voucher = await createVoucher(
        {
          owner: recipientAddress,
          tokenId: tokenId,
          amount: amount,
          uri: uri,
        },
        domain,
        owner
      );

      // The recipient address

      await expect(
        lazyMint1155.connect(owner).mintIfNotExists(voucher, recipientAddress)
      ).to.be.revertedWith("Invalid signature");
    });
    it("Should panic with invalid signature structure", async function () {
      const { lazyMint1155, owner, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const abiEncoded = AbiCoder.defaultAbiCoder().encode(
        ["uint256", "string", "uint256"],
        [tokenId, "hello", 1]
      );
      const hash = keccak256(abiEncoded);

      // Convert to Ethereum signed message format (same as toEthSignedMessageHash in Solidity)
      const messageHashBytes = getBytes(hash);
      const signature = await signer.signMessage(messageHashBytes);
      // The recipient address
      const recipientAddress = "0x476346a4510AeC7F469716935BF613656b4c22BD";

      const voucher = {
        owner: recipientAddress,
        tokenId: tokenId,
        amount: amount,
        uri: uri,
        signature: signature,
      };
      await expect(
        lazyMint1155.connect(owner).mintIfNotExists(voucher, recipientAddress)
      ).to.be.revertedWith("Invalid signature");
    });
  });

  describe("Global approvals", function () {
    it("Should allow the owner to set global approval for an operator", async function () {
      const { lazyMint1155, owner, account1 } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const operator = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      await lazyMint1155.connect(owner).setGlobalApproval(operator, true);

      const holder = account1;
      expect(await lazyMint1155.isApprovedForAll(holder.address, operator)).to
        .be.true;
      expect(await lazyMint1155.isGlobalApprover(operator)).to.be.true;
    });

    it("Should allow the owner to revoke global approval for an operator", async function () {
      const { lazyMint1155, owner, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const operator = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      await lazyMint1155.connect(owner).setGlobalApproval(operator, false);

      expect(await lazyMint1155.isApprovedForAll(owner.address, operator)).to.be
        .false;
      expect(await lazyMint1155.isGlobalApprover(operator)).to.be.false;
    });

    it("Should not allow non-owner to set global approval", async function () {
      const { lazyMint1155, signer } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const operator = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      await expect(
        lazyMint1155.connect(signer).setGlobalApproval(operator, true)
      ).to.be.revertedWithCustomError(
        lazyMint1155,
        "OwnableUnauthorizedAccount"
      );
    });

    it("Should allow holder opt out of global approval", async function () {
      const { lazyMint1155, owner, account1 } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const operator = "0x476346a4510AeC7F469716935BF613656b4c22BD";
      await lazyMint1155.connect(owner).setGlobalApproval(operator, true);
      expect(await lazyMint1155.isGlobalApprover(operator)).to.be.true;

      const holder = account1;
      expect(await lazyMint1155.isApprovedForAll(holder.address, operator)).to
        .be.true;

      await lazyMint1155.connect(holder).setGlobalApprovalOptOut(true);

      expect(await lazyMint1155.isApprovedForAll(holder.address, operator)).to
        .be.false;
    });

    it("Should allow holder toggle global approval", async function () {
      const { lazyMint1155, owner, account1 } = await loadFixture(
        deployLazyMint1155Fixture
      );

      const holder = account1;

      await lazyMint1155.connect(holder).setGlobalApprovalOptOut(true);

      expect(await lazyMint1155.hasOptedOutOfGlobalApproval(holder.address)).to
        .be.true;

      await lazyMint1155.connect(holder).setGlobalApprovalOptOut(false);

      expect(await lazyMint1155.hasOptedOutOfGlobalApproval(holder.address)).to
        .be.false;
    });

    it("Should respect standard ERC1155 approvals independently of global approvals", async function () {
      const { lazyMint1155, account1 } = await loadFixture(
        deployLazyMint1155Fixture
      );

      // Create a regular operator (not a global approver)
      const regularOperator = hre.ethers.Wallet.createRandom().address;

      // Verify operator is not a global approver and not approved yet
      expect(await lazyMint1155.isGlobalApprover(regularOperator)).to.be.false;
      expect(
        await lazyMint1155.isApprovedForAll(account1.address, regularOperator)
      ).to.be.false;

      // Set standard ERC1155 approval
      await lazyMint1155
        .connect(account1)
        .setApprovalForAll(regularOperator, true);

      // Check that the operator is now approved through standard mechanism
      expect(
        await lazyMint1155.isApprovedForAll(account1.address, regularOperator)
      ).to.be.true;
    });
  });

  // Additional tests for minting, transferring, etc. can be added here
});
