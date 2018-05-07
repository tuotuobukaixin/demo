#ifndef __PAAS_AES_CRYPTO_H__
#define __PAAS_AES_CRYPTO_H__

#define PRIMARY_KEYSTORE_FILE  "primary_store.txt"
#define STANDBY_KEYSTORE_FILE  "standby_store.txt"

#define PRIMARY_KMC_CFG_FILE  "primary_kmc_cfg.txt"
#define STANDBY_KMC_CFG_FILE  "standby_kmc_cfg.txt"


#define PAAS_AES_CRYPTO_LOG   "paas_crypto_log.log"
#define WSEC_AES_CRYPTO_LOG   "wsec_crypto_log.log"

#ifdef __cplusplus
extern "C"{
#endif /* __cplusplus */

/**  返回值 */
enum
{
    RET_SUCCESS         = 0,  /**< 成功 */
    RET_INVALID_PARAM   = 1,  /**< 入参错误 */
    RET_NORMAL_FAILURE  = 2,  /**< 内部一般异常 */
};

/* 加解密初始化 */
int aesInit();

/* 注册工作密钥 */
int registerWorkingKey(int iDomainId, int iKeyId, char *pPlainTexKey, int iKeyLen);

/* 设置工作密钥无效 */
int setKeyInvalid(int iDomainId, int iKeyId);

/* aes加密 */
char *aesEncrypt(int iDomainId, char *pcText);

/* aes解密 */
char *aesDecrypt(int iDomainId, char *pcHexEncData);

/* aes file encrypt */
int aesFileEncrypt(int iDomainId, const char *pszPlainFile, const char *pszCipherFile);

/* aes file decrypt */
int aesFileDecrypt(int iDomainId, const char *pszCipherFile, const char *pszPlainFile);

#ifdef __cplusplus
}
#endif /* __cplusplus */


#endif /* __PAAS_AES_CRYPTO_H__ */
