pragma solidity ^0.8.0;

import "hardhat/console.sol";

contract Decoder {
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

    function decodeBlockAndVerifySignature(bytes calldata rawBlock) external view returns (Block memory _block, bool success, string memory reason) {
        (bytes32 hash, bytes memory signature, bytes calldata rawBlockData) = decodeBlock(rawBlock);
        BlockData memory blockData = decodeBlockData(rawBlockData);
        _block.Hash = hash;
        _block.Signature = signature;
        _block.blockData = blockData;
        (success, reason) = verifySignature(rawBlockData, _block, signature);
    }

    function decodeBlockData(bytes calldata rawBlockData) public pure returns (BlockData memory blockData) {
        blockData.Data = rawBlockData[32 * 4 :];
        bytes memory tempWithoutData = rawBlockData[: 32 * 4];
        bytes32 _previousBlockHash;
        uint64 _blockNumber;
        uint64 _nonce;
        address _producerAddress;
        assembly {
            _previousBlockHash := mload(add(tempWithoutData, 32))
            _blockNumber := mload(add(tempWithoutData, 64))
            _nonce := mload(add(tempWithoutData, 96))
            _producerAddress := mload(add(tempWithoutData, 128))
        }
        blockData.PreviousBlockHash = _previousBlockHash;
        blockData.BlockNumber = _blockNumber;
        blockData.Nonce = _nonce;
        blockData.ProducerAddress = _producerAddress;
    }


    function decodeBlock(bytes calldata rawBlock) public pure returns (bytes32 hash, bytes memory signature, bytes calldata rawBlockData) {
        bytes memory tempHash = rawBlock[: 32];
        assembly {
            hash := mload(add(tempHash, 32))
        }
        signature = rawBlock[32 : 32 + 65];
        rawBlockData = rawBlock[32 + 65 :];
    }

    function verifySignature(bytes memory rawBlockData, Block memory block_, bytes memory signature) public view returns (bool, string memory) {
        bytes32 hash = keccak256(rawBlockData);
        (bytes32 r, bytes32 s, uint8 v, bool success) = parseSignature(signature);
        if (!success) {
            return (false, "decode signature failed");
        }
        address gotAddress = ecrecover(hash, v, r, s);
        if (gotAddress != block_.blockData.ProducerAddress) {
            return (false, "verify signature failed");
        }
        return (true, "");
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
