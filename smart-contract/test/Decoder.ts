import {ethers} from "hardhat";
import {assert} from "chai";

describe("Decoder", function () {
    let decoderContract: any;
    const blockDataSuccess = "0xd32004946b5de80a6790d9a4f3178b3f5f7a6d13f4f3c81bd9deaf508c2fdd061d0db15254df5de2c6194bfe436b0f13dd6aca5c0925846c76685a1597a0101c19c02342e228aab679b108fd83f61042c2b696309ecb714420e9636dfd35801d0101235398cd98047631439b784fbc186ac4fde6ec4ef5825cf999cb2b5870605500000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000047728aa300000000000000000000000071562b71999873db5b286df957af199ec94617f722736f6d656461746122"
    const blockDataFail = "0xd32004946b5de80a6780d9a4f3178b3f5f7a6d13f4f3c81bd9deaf508c2fdd061d0db15254df5de2c6194bfe436b0f13dd6aca5c0925846c76685a1597a0131c19c02342e228aab679b108fd83f61042c2b696309ecb714420e9636dfd15801d0101235398cd98047631439b784fbc186ac4fde6ec4ef5825cf999cb2b5870605500000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000047728aa300000000000000000000000071562b71999873db5b286df957af199ec94617f722736f6d656461746122"
    beforeEach(async function () {
        const decoder = await ethers.getContractFactory('Decoder')
        decoderContract = await decoder.deploy();
    })
    it("Should decode successfully", async function () {
        const result = await decoderContract.decodeBlockAndVerifySignature(blockDataSuccess)
        const blockData = result[0]
        assert.equal(result[1], true)
        assert.equal(blockData[0], "0xd32004946b5de80a6790d9a4f3178b3f5f7a6d13f4f3c81bd9deaf508c2fdd06")
        assert.equal(blockData[1], "0x1d0db15254df5de2c6194bfe436b0f13dd6aca5c0925846c76685a1597a0101c19c02342e228aab679b108fd83f61042c2b696309ecb714420e9636dfd35801d01")
        assert.equal(blockData[2][0], "0x01235398cd98047631439b784fbc186ac4fde6ec4ef5825cf999cb2b58706055")
        assert.equal(blockData[2][1], 2)
        assert.equal(blockData[2][2], 1198688931)
        assert.equal(blockData[2][3].toLowerCase(), "0x71562b71999873db5b286df957af199ec94617f7")
        assert.equal(blockData[2][4], "0x22736f6d656461746122")
    });
    it("Should decode fail", async function () {
        const result = await decoderContract.decodeBlockAndVerifySignature(blockDataFail)
        assert.equal(result[1], false)
        assert.equal(result[2], 'verify signature failed')
    });
});
