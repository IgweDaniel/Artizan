// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/utils/cryptography/MessageHashUtils.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/cryptography/EIP712.sol";

contract LazyMint1155 is ERC1155, Ownable, ReentrancyGuard, EIP712 {
    string private constant SIGNING_DOMAIN = "LazyMint1155";
    string private constant SIGNATURE_VERSION = "1";

    using ECDSA for bytes32;

    struct Voucher {
        address owner; 
        uint256 tokenId;
        uint256 amount;
        string uri;
        bytes signature; // signed by your backend or platform wallet
    }

    mapping(uint256 => bool) private _mintedTokens;

    mapping(address => bool) private _globalApprovers;

    mapping(address => bool) private _globalApprovalOptOut;

    address public signer;

    event GlobalApprovalOptOut(address indexed holder, bool optedOut);
    event GlobalApprovalSet(address indexed operator, bool approved);

    constructor(address _signer) ERC1155("") Ownable(msg.sender) EIP712(SIGNING_DOMAIN, SIGNATURE_VERSION) {
        signer = _signer;
    }

    function setSigner(address _signer) external onlyOwner {
        signer = _signer;
    }

    function _hash(Voucher calldata v) internal view returns (bytes32) {
        return _hashTypedDataV4(keccak256(abi.encode(
            keccak256("Voucher(address owner,uint256 tokenId,uint256 amount,string uri)"),
            v.owner,
            v.tokenId,
            v.amount,
            keccak256(bytes(v.uri))
        )));
    }

    function mintIfNotExists(Voucher calldata v, address to) public nonReentrant{
        if (_mintedTokens[v.tokenId]) {
            return; // Already minted, do nothing
        }
        // Require that the recipient matches the owner in the voucher
        require(to == v.owner, "Voucher owner mismatch");

        // Rebuild hash and recover signer
        bytes32 hash = _hash(v);
        address recovered = ECDSA.recover(hash, v.signature);
        require(recovered == signer, "Invalid signature");

        _mintedTokens[v.tokenId] = true;
        _mint(to, v.tokenId, v.amount, "");
        emit URI(v.uri, v.tokenId);
    }

    // Public function to check if a token has been minted
    function isTokenMinted(uint256 tokenId) public view returns (bool) {
        return _mintedTokens[tokenId];
    }
  
    // Set or revoke global approval for an operator
    function setGlobalApproval(address operator, bool approved) external onlyOwner {
        _globalApprovers[operator] = approved;
        emit GlobalApprovalSet(operator, approved);
    }

    // Allow users to opt out of global approval mechanism
    function setGlobalApprovalOptOut(bool optOut) external  {
        _globalApprovalOptOut[msg.sender] = optOut;
        emit GlobalApprovalOptOut(msg.sender, optOut);
    }

    // Check if a user has opted out of global approval
    function hasOptedOutOfGlobalApproval(address user) public view returns (bool) {
        return _globalApprovalOptOut[user];
    }

    // Check if an address is a global approver
    function isGlobalApprover(address operator) public view returns (bool) {
        return _globalApprovers[operator];
    }

    // Override isApprovedForAll to also check global approvers, respecting opt-outs
    function isApprovedForAll(address account, address operator) public view override returns (bool) {
        // If the account has opted out of global approvals, only check standard approvals
        if (_globalApprovalOptOut[account]) {
            return super.isApprovedForAll(account, operator);
        }
        
        // Otherwise check both standard approvals and global approver status
        return super.isApprovedForAll(account, operator) || _globalApprovers[operator];
    }
}
