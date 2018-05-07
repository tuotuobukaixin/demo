#ifndef __PAAS_AES_CRYPTO_H__
#define __PAAS_AES_CRYPTO_H__

#define PAAS_CRYPTO_PATH       "/var/paas"
#define PRIMARY_KEYSTORE_FILE  "/var/paas/primary_store.txt"
#define STANDBY_KEYSTORE_FILE  "/var/paas/standby_store.txt"

#define PRIMARY_KMC_CFG_FILE  "/var/paas/primary_kmc_cfg.txt"
#define STANDBY_KMC_CFG_FILE  "/var/paas/standby_kmc_cfg.txt"


#define PAAS_AES_CRYPTO_LOG   "/var/paas/paas_crypto_log.log"
#define WSEC_AES_CRYPTO_LOG   "/var/paas/wsec_crypto_log.log"

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

/**
************************************************************
 *@ingroup
 *@brief 加解密初始化
 *
 *
 *
 *@retval int          初始化返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int aesInit();

/**
************************************************************
 *@ingroup
 *@brief 注册工作密钥
 *
 *@param iDomainId       域ID
 *@param iKeyId          keyID
 *@param pPlainTexKey    key
 *@param iKeyLen         key长度
 *
 *@retval int          注册key返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int registerWorkingKey(int iDomainId, int iKeyId, char *pPlainTexKey, int iKeyLen);

/**
************************************************************
 *@ingroup
 *@brief 设置密钥状态为inactive，此时密钥只能做解密，不能继续加密
 *
 *@param iDomainId       域ID
 *@param iKeyId          keyID
 *
 *@retval int          设置key返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int setKeyInvalid(int iDomainId, int iKeyId);

/**
************************************************************
 *@ingroup
 *@brief aes加密
 *
 *@param iDomainId       域ID
 *@param pcText          明文
 *
 *@retval char*          密文
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
char *aesEncrypt(int iDomainId, char *pcText);

/**
************************************************************
 *@ingroup
 *@brief aes解密
 *
 *@param iDomainId       域ID
 *@param pcText          密文
 *
 *@retval char*          明文
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
char *aesDecrypt(int iDomainId, char *pcHexEncData);

/**
************************************************************
 *@ingroup
 *@brief aes加密文件
 *
 *@param iDomainId       域ID
 *@param pszPlainFile    原始文件名，需要指定路径
 *@param pszPlainFile    加密文件名，需要指定路径
 *
 *@retval int            加密成功失败返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int aesFileEncrypt(int iDomainId, const char *pszPlainFile, const char *pszCipherFile);

/**
************************************************************
 *@ingroup
 *@brief aes解密文件
 *
 *@param iDomainId       域ID
 *@param pszPlainFile    加密文件名，需要指定路径
 *@param pszPlainFile    原始文件名，需要指定路径
 *
 *@retval int            加密成功失败返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int aesFileDecrypt(int iDomainId, const char *pszCipherFile, const char *pszPlainFile);

#ifdef __cplusplus
}
#endif /* __cplusplus */


#endif /* __PAAS_AES_CRYPTO_H__ */
