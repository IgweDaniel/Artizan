// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;
// import "seaport-types/interfaces/ZoneInterface.sol";
import {ZoneInterface} from "seaport-types/src/interfaces/ZoneInterface.sol";
import {ZoneParameters,Schema} from "seaport-types/src/lib/ConsiderationStructs.sol";

import { ERC165 } from "@openzeppelin/contracts/utils/introspection/ERC165.sol";

import {LazyMint1155} from "./LazyMint1155.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract LazyMintZone is ERC165,ZoneInterface,Ownable {
    LazyMint1155 public nft;

    constructor(address _signer, address _nft)  Ownable(_signer) {
        nft = LazyMint1155(_nft);
    }

    function validateOrder(
        ZoneParameters calldata
    ) external pure returns (bytes4) {
        return this.validateOrder.selector;
    }
    function authorizeOrder(
        ZoneParameters calldata zoneParameters
    ) external override returns (bytes4) {
        
        LazyMint1155.Voucher memory voucher = abi.decode(
            zoneParameters.extraData,
            (LazyMint1155.Voucher)
        );
        
        nft.mintIfNotExists(voucher, zoneParameters.offerer);
        return this.authorizeOrder.selector;
    }
 

    function getSeaportMetadata()
        external
        pure
        override
        returns (
            string memory name,
            Schema[] memory schemas // map to Seaport Improvement Proposal IDs
        )
    {
        schemas = new Schema[](1);
        schemas[0].id = 3003;
        schemas[0].metadata = new bytes(0);

        return ("ArtiartZone", schemas);
    }

     function supportsInterface(
        bytes4 interfaceId
    ) public view override(ERC165, ZoneInterface) returns (bool) {
        return
            interfaceId == type(ZoneInterface).interfaceId ||
            super.supportsInterface(interfaceId);
    }

    function setNftAddress(address _nft) external onlyOwner {
        require(_nft != address(0), "Invalid NFT address");
        nft = LazyMint1155(_nft);
    }
}
