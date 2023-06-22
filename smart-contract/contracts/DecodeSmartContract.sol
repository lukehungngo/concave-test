pragma solidity ^0.8.0;

contract DecodeSmartContract {
    struct BlockData {
        bytes32 PreviousBlockHash;
        uint64 BlockNumber;
        uint64 Nonce;
        address ProducerAddress;
        bytes Data;
    }

    struct Block {
        bytes32 Hash;
        bytes Signature;
        BlockData blockData;
    }
    // Test: 0xb6dafeef592228dadfc0e3e53d50d2e971652dac745b2a5ffcbb602365a5c2af7aa909a341af5f55173ee086773e5e1a2f26ccb2bde151fe2aedec3fdd22530066a523cd06f73ca7c0dae4476116024fa1d6a694ac6378e3fa74b7a1b3b1108200b435c006d8008ecdaa02f09d8bc55089b9d2101c190d2af73fd7521ecb6b1387000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000036d68e9b2d43a89d00000000000000000000000071562b71999873db5b286df957af199ec94617f7226162636422
    function decodeBlockAndVerifySignature(bytes calldata rawBlock) external pure returns(Block memory _block) {
        (bytes32 hash, bytes memory signature, bytes calldata rawBlockData) = decodeBlock(rawBlock);
        BlockData memory blockData = decodeBlockData(rawBlockData);
        _block.Hash = hash;
        _block.Signature = signature;
        _block.blockData = blockData;
        bool success = verifySignature(rawBlock, _block, signature);
        require(success, "verify signature failed");
    }

    function decodeBlockData(bytes calldata rawBlockData) public pure returns(BlockData memory blockData) {
        blockData.Data = rawBlockData[32*4:];
        bytes memory tempWithoutData = rawBlockData[:32*4];
        bytes32 _previousBlockHash;
        uint64 _blockNumber;
        uint64 _nonce;
        address _producerAddress;
        assembly {
            _previousBlockHash := mload(add(tempWithoutData,32))
            _blockNumber := mload(add(tempWithoutData,64))
            _nonce := mload(add(tempWithoutData,96))
            _producerAddress := mload(add(tempWithoutData,128))
        }
        blockData.PreviousBlockHash = _previousBlockHash;
        blockData.BlockNumber = _blockNumber;
        blockData.Nonce = _nonce;
        blockData.ProducerAddress = _producerAddress;
    }


    function decodeBlock(bytes calldata rawBlock) public pure returns(bytes32 hash, bytes memory signature, bytes calldata rawBlockData) {
        bytes memory tempHash = rawBlock[:32];
        assembly {
            hash := mload(add(tempHash,32))
        }
        signature = rawBlock[32:32+65];
        rawBlockData = rawBlock[32+65:];
    }

    function verifySignature(bytes memory rawBlock, Block memory block_, bytes memory signature) public pure returns(bool) {
        bytes32 hash = keccak256(rawBlock);
        (bytes32 r, bytes32 s, uint8 v, bool success) = parseSignature(signature);
        require(success, "decode signature failed");
        address gotAddress = ecrecover(hash, v, r, s);
        return (gotAddress != block_.blockData.ProducerAddress);
    }

    function parseSignature(bytes memory signature) public pure returns (bytes32 r, bytes32 s, uint8 v, bool success) {
        require(signature.length == 65, "Invalid signature length");

        assembly {
        // Load the first 32 bytes (r) from the signature
        // Skip first 32 bytes in the dynamic array, because dynamic array use first 32 bytes to store length of array
            r := mload(add(signature, 32))
        // Load the next 32 bytes (s) from the signature
            s := mload(add(signature, 64))
        // Load the last byte (v) from the signature
        // Because mload will load 32 bytes, but we only need one last byte, so we extract the last bytes only after mload
            v := byte(0, mload(add(signature, 96)))
        }

        if (v < 27) {
            v += 27;
        }

        success = true;
        if (v != 27 && v != 28) {
            success = false;
        }

    }
}