require 'dl'
require 'dl/import'

class Crypto
  module Libaes_crypto
    extend DL::Importer
    dlload '/usr/lib/libaes_crypto_dyn.so'
    extern 'int aesInit()'
    extern 'int registerWorkingKey(int, int, char *, int *)'
    extern 'char* aesEncrypt(int, char *)'
    extern 'char* aesDecrypt(int, char *)'
    extern 'int   aesFileEncrypt(int, char *, char *)'
    extern 'int   aesFileDecrypt(int, char *, char *)'
  end 

  def aes_init() 
    ret = Libaes_crypto.aesInit()
    return ret
  end
  
  def aes_register_key(domainId, keyId, plainKey, keyLen) 
    ret = 0
    ret = Libaes_crypto.RegisterWorkingKey(domainId, keyId, plainKey, keyLen)
    return ret
  end

  def aes_encrypt(domainId, plaintext)
    encryptData_c = Libaes_crypto.aesEncrypt(domainId, plaintext)
    encryptData = encryptData_c.to_s
    DL.free(encryptData_c)
    return encryptData
  end

  def aes_decrypt(domainId, encryptData)
    plaintext_c = Libaes_crypto.aesDecrypt(domainId, encryptData)
    plaintext = plaintext_c.to_s
    DL.free(plaintext_c)
    return plaintext
  end
  
  def aes_file_encrypt(domainId, plainFile, encryptFile)
    ret = Libaes_crypto.aesFileEncrypt(domainId, plainFile, encryptFile)
    return ret
  end
  
  def aes_file_decrypt(domainId, encryptFile, plainFile)
    ret = Libaes_crypto.aesFileDecrypt(domainId, encryptFile, plainFile)
    return ret
  end
end