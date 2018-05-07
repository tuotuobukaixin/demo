from ctypes import *

crypto = CDLL("/usr/lib/libaes_crypto_dyn.so")

crypto.aesFree.argtypes=[c_void_p]
crypto.aesEncrypt.argtypes=[c_int, c_char_p]
crypto.aesDecrypt.argtypes=[c_int, c_char_p]
crypto.aesFileEncrypt.argtypes=[c_int, c_char_p, c_char_p]
crypto.aesFileDecrypt.argtypes=[c_int, c_char_p, c_char_p]
crypto.aesSetRootKeyExpiredDays.argtypes=[c_int]
crypto.aesInit.restype = c_int
crypto.aesFree.restype = c_int
crypto.aesEncrypt.restype = c_void_p
crypto.aesDecrypt.restype = c_void_p
crypto.aesFileEncrypt.restype = c_int
crypto.aesFileDecrypt.restype = c_int
crypto.aesSetRootKeyExpiredDays.restype = c_int
class Crypto:
    def aes_init(self):
        return crypto.aesInit()

    def aes_encrypt(self, domianId, plainText):
        encryptData_c = crypto.aesEncrypt(domianId, plainText)
        encryptData = cast(encryptData_c, c_char_p).value
        crypto.aesFree(encryptData_c)
        return encryptData

    def aes_decrypt(self, domianId, encryptData):
        plainText_c = crypto.aesDecrypt(domianId, encryptData)
        plainText =  cast(plainText_c, c_char_p).value
        crypto.aesFree(plainText_c)
        return plainText
            
    def aes_file_encrypt(self, domianId, plainFile, encFile):
        ret = crypto.aesFileEncrypt(domianId, plainFile, encFile)
        return ret

    def aes_file_decrypt(self, domianId, encFile, plainFile):
        ret = crypto.aesFileDecrypt(domianId, encFile, plainFile)
        return ret

    def aes_update_rootkey(self):
        crypto.aesUpdateRootKey()
        return

    def aes_set_rootkey_expired_days(self, expiredDays):
        ret = crypto.aesSetRootKeyExpiredDays(expiredDays)
        return ret
