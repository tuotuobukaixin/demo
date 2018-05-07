#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <time.h>
#include <dlfcn.h>
#include <sys/stat.h>
#include "aes_crypto.h"
#include "WSEC_Itf.h"
#include "WSEC_ErrorCode.h"
#include "SDP_Itf.h"
#include "KMC_Itf.h"
#include "WSEC_Type.h"
#include "securec.h"
#include "securectype.h"

#define BUFFER_LEN  2048
#define MAX_LOG_FILE_SIZE (1024*1024*1024)

#ifdef __cplusplus
extern "C"{
#endif /* __cplusplus */

/**
************************************************************
 *@ingroup
 *@brief 查询文件大小
 *
 *
 *@param nLevel           日志级别
 *@param pszModuleName    模块名
 *@param pszOccurFileName 文件名
 *@param nOccurLine       行号
 *@param pszLog           日志内容
 *
 *@retval void
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
unsigned long getFileSize(char *pcFileName)
{  
    unsigned long filesize = 0;
    struct stat statbuff;
    if(stat(pcFileName, &statbuff) > 0){
        filesize = statbuff.st_size;
    }
    return filesize;
}

/**
************************************************************
 *@ingroup
 *@brief 注册给WSEC的写日志函数
 *
 *
 *@param nLevel           日志级别
 *@param pszModuleName    模块名
 *@param pszOccurFileName 文件名
 *@param nOccurLine       行号
 *@param pszLog           日志内容
 *
 *@retval void
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
void wsecCryptoWriteLog(int nLevel, const char* pszModuleName, const char* pszOccurFileName, int nOccurLine, const char* pszLog)
{
    int iRet = 0;
    int   ibufLen;
    char  szbuffer[BUFFER_LEN] = {0};
    unsigned long ulFileSize = 0;
    FILE *pfCryptoLog = NULL;
    ulFileSize = getFileSize(WSEC_AES_CRYPTO_LOG);
    if (ulFileSize > MAX_LOG_FILE_SIZE)
    {
        iRet = remove(WSEC_AES_CRYPTO_LOG);
        if (iRet != 0)
        {
            printf("remove log file failed!\n");
        }
    }
    pfCryptoLog = fopen(WSEC_AES_CRYPTO_LOG, "a+");
    if (NULL == pfCryptoLog)
    {
        printf("open file wsec log error!\n");
        return;
    }

    ibufLen = snprintf_s(szbuffer, BUFFER_LEN-1, BUFFER_LEN-1, "%s %d %s", pszModuleName, nOccurLine, pszLog);
    if (ibufLen <= 0)
    {
        printf("snprintf string error!\n");
        fclose(pfCryptoLog);
        return;
    }

#if 0
    printf(szbuffer, "%s %d %s", pszModuleName, nOccurLine, pszLog);
#else
    iRet = fwrite(szbuffer, ibufLen, 1, pfCryptoLog);
    if (iRet != 1)
    {
        printf("write file failed!\n");
    }
#endif
    fclose(pfCryptoLog);
    return;
}

/**
************************************************************
 *@ingroup
 *@brief PAAS内部的加解密写日志函数
 *
 *
 *@param pcFormat         格式化字符串
 *@param ...              参数
 *
 *@retval void
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
void paasCryptoWriteLog(char *pcFormat, ...)
{
    int iRet = 0;
    va_list  argument;
    int      iTimeBufLen;
    int      iLogBufLen;
    char     szbuffer[BUFFER_LEN] = {0};
    char     szTimeBuf[BUFFER_LEN] = {0};
    time_t rawtime;
    struct tm * timeinfo;
    unsigned long ulFileSize = 0;
    FILE *pfCryptoLog = NULL;

    ulFileSize = getFileSize(WSEC_AES_CRYPTO_LOG);
    if (ulFileSize > MAX_LOG_FILE_SIZE)
    {
        iRet = remove(PAAS_AES_CRYPTO_LOG);
        if (iRet != 0)
        {
            paasCryptoWriteLog("remove crypto log file failed!\n");
        }
    }

    time(&rawtime);
    timeinfo = localtime(&rawtime);
    if (NULL == timeinfo)
    {
        printf("get local time error!\n");
        return;
    }
    iTimeBufLen = snprintf_s(szTimeBuf, BUFFER_LEN-1, BUFFER_LEN-1, "%4d-%02d-%02d %02d:%02d:%02d [paas] ", 1900 + timeinfo->tm_year,
        1 + timeinfo->tm_mon, timeinfo->tm_mday, timeinfo->tm_hour, timeinfo->tm_min, timeinfo->tm_sec);
    if (iTimeBufLen <= 0)
    {
        printf("snprintf time error!\n");
        return;
    }
    
    va_start(argument, pcFormat);
    iLogBufLen = vsnprintf_s(szbuffer, BUFFER_LEN-1, BUFFER_LEN-1, (const char*)pcFormat, argument);
    va_end(argument);
    if (iLogBufLen <= 0)
    {
        printf("vsnprintf string error!\n");
        return;
    }
#if 1
    pfCryptoLog = fopen(PAAS_AES_CRYPTO_LOG, "a+");
    if (NULL == pfCryptoLog)
    {
        printf("open file paas log error!\n");
        return;
    }
    iRet = fwrite(szTimeBuf, iTimeBufLen, 1, pfCryptoLog);
    if (iRet != 1)
    {
        paasCryptoWriteLog("write crypto file failed!\n");
    }
    iRet = fwrite(szbuffer, iLogBufLen, 1, pfCryptoLog);
    if (iRet != 1)
    {
        paasCryptoWriteLog("write crypto file failed!\n");
    }
    fclose(pfCryptoLog);
#else
    printf("%s", szbuffer);
#endif
    return;
}

/**
************************************************************
 *@ingroup
 *@brief 注册给WSEC密钥过期通知函数
 *
 *@param eNtfCode         通知类型
 *@param pData            通知内容
 *@param nDataSize        内容大小
 *
 *@retval void
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
void wsecCryptoNotify(WSEC_NTF_CODE_ENUM eNtfCode, const WSEC_VOID* pData, WSEC_SIZE_T nDataSize)
{
    if (pData == NULL)
    {
        return;
    }
    KMC_MK_EXPIRE_NTF_STRU *pstKmcMkExpireNtf = NULL;
    KMC_MK_CHANGE_NTF_STRU *pstKmcMkChangeNtf = NULL;
    KMC_RK_EXPIRE_NTF_STRU *pstKmcRkExpireNtf = NULL;
    
    switch (eNtfCode)
    {
        case WSEC_KMC_NTF_MK_EXPIRE:
        {
            pstKmcMkExpireNtf = (KMC_MK_EXPIRE_NTF_STRU *)pData;
            paasCryptoWriteLog("Register key domain %d keyid %d have %d days to expired!\n",
                pstKmcMkExpireNtf->stMkInfo.ulDomainId, pstKmcMkExpireNtf->stMkInfo.ulKeyId,
                pstKmcMkExpireNtf->nRemainDays);
            break;
        }

        case WSEC_KMC_NTF_MK_CHANGED:
        {
            pstKmcMkChangeNtf = (KMC_MK_CHANGE_NTF_STRU *)pData;
            paasCryptoWriteLog("Register key domain %d keyid %d status changed to %d!\n",
                pstKmcMkChangeNtf->stMkInfo.ulDomainId, pstKmcMkChangeNtf->stMkInfo.ulKeyId,
                pstKmcMkChangeNtf->eType);
            break;
        }

        case WSEC_KMC_NTF_RK_EXPIRE:
        {
            pstKmcRkExpireNtf = (KMC_RK_EXPIRE_NTF_STRU *)pData;
            if (pstKmcRkExpireNtf->nRemainDays == 0)
            {
                paasCryptoWriteLog("Root Key expired! Old Key is created in %d-%02d-%02d %02d:%02d:%02d (UTC), Expired in %d-%02d-%02d %02d:%02d:%02d (UTC)\n",
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.uwYear,
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.ucMonth,
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.ucDate,
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.ucHour,
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.ucMinute,
                    pstKmcRkExpireNtf->stRkInfo.stRkCreateTimeUtc.ucSecond,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.uwYear,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.ucMonth,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.ucDate,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.ucHour,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.ucMinute,
                    pstKmcRkExpireNtf->stRkInfo.stRkExpiredTimeUtc.ucSecond);
            }
            else
            {
                paasCryptoWriteLog("Root key will expired in %d days!\n",  pstKmcRkExpireNtf->nRemainDays);
            }
            break;
        }

        default:
        {
            paasCryptoWriteLog("Have a notify event, notify code=%d!\n",
                    eNtfCode);
        }
    }
    return;
}

/**
************************************************************
 *@ingroup
 *@brief 注册给WSEC密钥过期事件处理函数
 *
 *
 *
 *@retval void
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
void wsecCryptoDoevent()
{
    return;
}

/**
************************************************************
 *@ingroup
 *@brief 字符串转16进制字符串
 *
 *@param pcInputStr       原始字符串
 *@param iStrLen          字符串长度
 *
 *@retval char*
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
char* str2hex(char* pcInputStr, int iStrLen)
{
    errno_t      iRet = 0;
    int          i = 0;
    unsigned int j;
    unsigned int a;
    const char *pcHex = "0123456789ABCDEF";
    char    *pchex = (char*) malloc(2*iStrLen + 1);
    if (NULL == pchex)
    {
        paasCryptoWriteLog("str2hex malloc mem failed!\n");
        return NULL;
    }
    iRet = memset_s(pchex, 2*iStrLen + 1, 0, 2*iStrLen + 1);
    if (iRet != EOK)
    {
        paasCryptoWriteLog("str2hex memset failed!\n");
        free(pchex);
        pchex = NULL;
        return NULL;        
    }
    for(j = 0; j < iStrLen; j++)
    {
        a =  (unsigned int) pcInputStr[j];
        pchex[i++] = pcHex[(a & 0xf0) >> 4];
        pchex[i++] = pcHex[(a & 0x0f)];
    }
    pchex[i] = '\0';
    return pchex;
}

/**
************************************************************
 *@ingroup
 *@brief 16进制字符转换成整数
 *
 *@param hex       16进制字符
 *
 *@retval int
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int convertHex2Int(char hex)
{
    if (hex >= '0' && hex <= '9')
    {
        return hex-'0';
    }
    else if (hex >= 'A' && hex <= 'F')
    {
        return hex-'A'+10;
    }
    else
    {
        paasCryptoWriteLog("The input char %c is invliad!\n", hex);
        return -1;
    }
}


/**
************************************************************
 *@ingroup
 *@brief 16进制字符串转字符串
 *
 *@param pcHexStr       16进制字符串
 *@param iStrLen        字符串长度
 *
 *@retval unsigned char*
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
unsigned char* hex2str(char* pcHexStr, int iStrLen)
{
    errno_t      iRet = 0;
    int          i = 0;
    unsigned int j;
    unsigned int high;
    unsigned int low;
    unsigned char *pucStr = (unsigned char*) malloc(iStrLen/2 + 1);
    if (NULL == pucStr)
    {
        paasCryptoWriteLog("hex2str malloc mem failed!\n");
        return NULL;
    }
    iRet = memset_s(pucStr, iStrLen/2 + 1, 0, iStrLen/2 + 1);
    if (iRet != EOK)
    {
        paasCryptoWriteLog("hex2str memset failed!\n");
        free(pucStr);
        pucStr = NULL;
        return NULL;        
    }
    for(j = 0; j < iStrLen; j = j + 2)
    {
        high = convertHex2Int(pcHexStr[j]);
        low = convertHex2Int(pcHexStr[j+1]);
        if (high > 15 || low > 15)
        {
            paasCryptoWriteLog("The input Encrypt Data is invalid!\n");
            free(pucStr);
            pucStr = NULL;
            return NULL;
        }
        pucStr[i++] = (high << 4) + low;
    }
    pucStr[i] = '\0';
    return pucStr;
}

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
int registerWorkingKey(int iDomainId, int iKeyId, char *pPlainTexKey, int iKeyLen)
{
    int iRet = 0;
    unsigned short usKeyType = KMC_KEY_TYPE_ENCRPT_INTEGRITY;

    iRet = KMC_RegisterMk(iDomainId, iKeyId, usKeyType, pPlainTexKey, iKeyLen);
    if (iRet != RET_SUCCESS && iRet != WSEC_ERR_KMC_REG_REPEAT_MK)
    {
        paasCryptoWriteLog("Register working key failed, ret=%d\n", iRet);
        return RET_NORMAL_FAILURE;
    }
    paasCryptoWriteLog("Register working key success\n");
    return RET_SUCCESS;
}

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
int setKeyInvalid(int iDomainId, int iKeyId)
{
    int iRet = 0;

    iRet = KMC_SetMkStatus(iDomainId, iKeyId, KMC_KEY_STATUS_INACTIVE);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("set key inactive failed, ret=%d\n", iRet);
        return RET_NORMAL_FAILURE;
    }
    paasCryptoWriteLog("set key inactive success\n");
    return RET_SUCCESS;
}

/**
************************************************************
 *@ingroup
 *@brief 检查加密所需目录是否存在，不存在则创建
 *
 *@param pcPath          目录
 *
 *
 *@retval int            目录处理返回码
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int checkDirectory(char *pcPath)
{
    FILE *fp = NULL;
    if (NULL == pcPath)
    {
        return RET_NORMAL_FAILURE;
    }
    fp = fopen(pcPath, "r");
    if (NULL == fp)
    {
        if(mkdir(pcPath, 0750) == -1)  
        {
            paasCryptoWriteLog("create directory failed\n");
            return RET_NORMAL_FAILURE;
        }
        return RET_SUCCESS;
    }
    fclose(fp);
    return RET_SUCCESS;
}

/**
************************************************************
 *@ingroup
 *@brief 释放内存
 *
 *
 *
 *@retval int          释放内存是否成功
  修改历史   :
  1.日    期   : 2015年04月28日
    作    者   : z00223295
    修改内容:
************************************************************/
int aesFree(void *ptr)
{
    if (ptr != NULL)
    {
        free(ptr);
        ptr = NULL;
        paasCryptoWriteLog("aes free ptr success\n");
    }
    return RET_SUCCESS;
}

/******************************************************************************
 *
 * @brief: update root key, it's a temporary function because crypto CBB have
 *         a bug, it can't update root key info in memory after call 
 *         KMC_UpdateRootKey, so I must call WSEC_Reset after it .
 *
 * @param: VOID
 *
 * @return:  0 OK, other failed
 *
 *****************************************************************************/
void aesUpdateRootKey()
{
    KMC_RK_ATTR_STRU        stKmcRkAttr = {0};

    (void)KMC_UpdateRootKey(NULL, 0);
    (void)WSEC_Reset();

    if (KMC_GetRootKeyInfo(&stKmcRkAttr) == RET_SUCCESS)
    {
        paasCryptoWriteLog("New root key is created! Created in  %d-%02d-%02d %02d:%02d:%02d (UTC), Expired in %d-%02d-%02d %02d:%02d:%02d (UTC)\n",
                stKmcRkAttr.stRkCreateTimeUtc.uwYear,
                stKmcRkAttr.stRkCreateTimeUtc.ucMonth,
                stKmcRkAttr.stRkCreateTimeUtc.ucDate,
                stKmcRkAttr.stRkCreateTimeUtc.ucHour,
                stKmcRkAttr.stRkCreateTimeUtc.ucMinute,
                stKmcRkAttr.stRkCreateTimeUtc.ucSecond,
                stKmcRkAttr.stRkExpiredTimeUtc.uwYear,
                stKmcRkAttr.stRkExpiredTimeUtc.ucMonth,
                stKmcRkAttr.stRkExpiredTimeUtc.ucDate,
                stKmcRkAttr.stRkExpiredTimeUtc.ucHour,
                stKmcRkAttr.stRkExpiredTimeUtc.ucMinute,
                stKmcRkAttr.stRkExpiredTimeUtc.ucSecond);
    }
}

/******************************************************************************
 *
 * @brief: Judge if root key is expired
 *
 * @param: VOID
 *
 * @return: Positive: not expire, minus and 0: expired
 *
 *****************************************************************************/
int judgeRootKeyExpired()
{
    int iRet = 0;
    time_t stTime = {0};
    struct tm *pstTimeDetail = NULL;
    KMC_RK_ATTR_STRU stRkAttr = {0};

    iRet = KMC_GetRootKeyInfo(&stRkAttr);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("judgeRootKeyExpired get root key info failed, return %d!\n", iRet);
        return RET_INVALID_PARAM;
    }

    time(&stTime);
    pstTimeDetail = gmtime(&stTime);
    if (pstTimeDetail == NULL)
    {
        paasCryptoWriteLog("System error! get localtime failed\n");
        return RET_INVALID_PARAM;
    }

    if (iRet = stRkAttr.stRkExpiredTimeUtc.uwYear - (pstTimeDetail->tm_year + 1900))    return iRet;
    if (iRet = stRkAttr.stRkExpiredTimeUtc.ucMonth - (pstTimeDetail->tm_mon + 1))       return iRet;
    if (iRet = stRkAttr.stRkExpiredTimeUtc.ucDate - pstTimeDetail->tm_mday)             return iRet;
    if (iRet = stRkAttr.stRkExpiredTimeUtc.ucHour - pstTimeDetail->tm_hour)             return iRet;
    if (iRet = stRkAttr.stRkExpiredTimeUtc.ucMinute - pstTimeDetail->tm_min)            return iRet;
    if (iRet = stRkAttr.stRkExpiredTimeUtc.ucSecond - pstTimeDetail->tm_sec)            return iRet;
    
    return iRet;
}

/******************************************************************************
 *
 * @brief: set root key expired days
 *
 * @param: int iDays (1 days to 10 years)
 *
 * @return: 0 OK, other failed
 *
 *****************************************************************************/
int aesSetRootKeyExpiredDays(int iDays)
{
    int iRet = 0;
    KMC_CFG_ROOT_KEY_STRU   stRootKeyInfo = {0};

    if (iDays <= 0 || iDays > 3650)
    {
        paasCryptoWriteLog("Root key expired days only between 1 to 3650, the input %d is out of range\n", iDays);
        return RET_INVALID_PARAM;
    } 

    iRet = KMC_GetRootKeyCfg(&stRootKeyInfo);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("Get root key info failed, ret = %d\n", iRet);
        return RET_NORMAL_FAILURE;
    }

    stRootKeyInfo.ulRootKeyLifeDays = iDays;
    iRet = KMC_SetRootKeyCfg(&stRootKeyInfo);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("Set root key info failed, ret = %d\n", iRet);
        return RET_NORMAL_FAILURE;
    }

    /* After setting new expire time , it must create a new root key for new config */
    aesUpdateRootKey();

    return RET_SUCCESS;
}

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
int aesInit()
{
    int iRet = 0;
    KMC_FILE_NAME_STRU  stKmcFileName = {0};
    WSEC_FP_CALLBACK_STRU stWsecFpCallback = {0};

    iRet = checkDirectory(PAAS_CRYPTO_PATH);
    if (iRet != RET_SUCCESS)
    {
        return RET_NORMAL_FAILURE;
    }

    stKmcFileName.pszKeyStoreFile[0] = PRIMARY_KEYSTORE_FILE;
    stKmcFileName.pszKeyStoreFile[1] = STANDBY_KEYSTORE_FILE;
    stKmcFileName.pszKmcCfgFile[0]   = PRIMARY_KMC_CFG_FILE;
    stKmcFileName.pszKmcCfgFile[1]   = STANDBY_KMC_CFG_FILE;

    stWsecFpCallback.stRelyApp.pfWriLog   = wsecCryptoWriteLog;
    stWsecFpCallback.stRelyApp.pfNotify   = wsecCryptoNotify;
    stWsecFpCallback.stRelyApp.pfDoEvents = wsecCryptoDoevent;
    iRet = WSEC_Initialize(&stKmcFileName, &stWsecFpCallback, NULL, NULL);
    if (iRet != RET_SUCCESS && iRet != WSEC_ERR_KMC_INI_MUL_CALL)
    {
        paasCryptoWriteLog("aes init failed ret=%d\n", iRet);
        return RET_NORMAL_FAILURE;
    }
    paasCryptoWriteLog("aes init success\n");
    return RET_SUCCESS;
}

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
char *aesEncrypt(int iDomainId, char *pcText)
{
    errno_t nRet = 0;

    int iRet        = 0;
    int iTextLen    = 0;
    int iEncDataLen = 0;
    
    char *pcEncData = NULL;
    char *pcHexData = NULL;
    if (NULL == pcText)
    {
        paasCryptoWriteLog("The input text is null!\n");
        return NULL;
    }
    
    if (judgeRootKeyExpired() <= 0)
    {
        aesUpdateRootKey();
    }

    iTextLen = strlen(pcText);
    iRet = SDP_GetCipherDataLen(iTextLen, &iEncDataLen);
    if (RET_SUCCESS != iRet || iEncDataLen <= 0)
    {
        paasCryptoWriteLog("Get cipher data Len failed, iRet=%d\n", iRet);
        return NULL;
    }
    pcEncData = (char *)malloc(iEncDataLen + 1);
    if (NULL == pcEncData)
    {
        paasCryptoWriteLog("aesEncrypt malloc mem failed!\n");
        return NULL;
    }
    nRet = memset_s(pcEncData, iEncDataLen + 1, 0, iEncDataLen + 1);
    if (nRet != EOK)
    {
        paasCryptoWriteLog("aesEncrypt memset failed!\n");
        free(pcEncData);
        pcEncData = NULL;
        return NULL;        
    }

    iRet = SDP_Encrypt(iDomainId, pcText, iTextLen, pcEncData, &iEncDataLen);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("Encrypt failed, iRet=%d\n", iRet);
        free(pcEncData);
        pcEncData = NULL;
        return NULL;
    }
    pcHexData = str2hex(pcEncData, iEncDataLen);
    free(pcEncData);
    pcEncData = NULL;
    paasCryptoWriteLog("Encrypt text success\n");
    return pcHexData;
}

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
char *aesDecrypt(int iDomainId, char *pcHexEncData)
{
    errno_t nRet = 0;

    int iRet         = 0;
    int iEncLen      = 0;
    int iTextLen     = 0; 
    int iHexEncLen   = 0;
    
    char *pcPlainText = NULL;
    unsigned char  *pucEncData   = NULL;

    if (NULL == pcHexEncData)
    {
        paasCryptoWriteLog("The input Encrypt Data is null!\n");
        return NULL;
    }

    if (judgeRootKeyExpired() <= 0)
    {
        aesUpdateRootKey();
    }

    iHexEncLen = strlen(pcHexEncData);
    if (iHexEncLen%2 != 0)
    {
        paasCryptoWriteLog("The input Encrypt Data not right\n");
    }
    iEncLen = iHexEncLen/2;

    pcPlainText = (char *)malloc(iEncLen + 1);
    if (NULL == pcPlainText)
    {
        paasCryptoWriteLog("aesDecrypt malloc mem failed!\n");
        return NULL;
    }

    nRet = memset_s(pcPlainText, iEncLen + 1, 0, iEncLen + 1);
    if (nRet != EOK)
    {
        paasCryptoWriteLog("aesDecrypt memset failed!\n");
        free(pcPlainText);
        pcPlainText = NULL;
        return NULL;       
    }
    iTextLen = iEncLen;
    pucEncData = hex2str(pcHexEncData, iHexEncLen);
    if (NULL == pucEncData)
    {
        free(pcPlainText);
        pcPlainText = NULL;
        return NULL;
    }
    iRet = SDP_Decrypt(iDomainId, pucEncData, iEncLen, pcPlainText, &iTextLen);
    free(pucEncData);
    pucEncData = NULL;
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("Decrypt failed, iRet=%d\n", iRet);
        free(pcPlainText);
        pcPlainText = NULL;
        return NULL;
    }
    paasCryptoWriteLog("Decrypt text success\n");
    return pcPlainText;
}

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
int aesFileEncrypt(int iDomainId, const char *pszPlainFile, const char *pszCipherFile)
{
    int iRet         = 0;
    if (NULL == pszPlainFile || NULL == pszCipherFile)
    {
        paasCryptoWriteLog("File encrypt param error!\n");
        return RET_INVALID_PARAM;
    }

    if (judgeRootKeyExpired() <= 0)
    {
        aesUpdateRootKey();
    }

    iRet = SDP_FileEncrypt(iDomainId, pszPlainFile, pszCipherFile, NULL, NULL);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("File encrypt failed, iRet=%d\n", iRet);
        return RET_NORMAL_FAILURE;
    }
    paasCryptoWriteLog("Encrypt file success\n");
    return RET_SUCCESS;
}

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
int aesFileDecrypt(int iDomainId, const char *pszCipherFile, const char *pszPlainFile)
{
    int iRet         = 0;
    if (NULL == pszCipherFile || NULL == pszPlainFile)
    {
        paasCryptoWriteLog("File decrypt param error!\n");
        return RET_INVALID_PARAM;
    }

    if (judgeRootKeyExpired() <= 0)
    {
        aesUpdateRootKey();
    }

    iRet = SDP_FileDecrypt(iDomainId, pszCipherFile, pszPlainFile, NULL, NULL);
    if (iRet != RET_SUCCESS)
    {
        paasCryptoWriteLog("File decrypt failed, iRet=%d\n", iRet);
        return RET_NORMAL_FAILURE;
    }
    paasCryptoWriteLog("Decrypt file success\n");
    return RET_SUCCESS;
}



#ifdef __cplusplus
}
#endif /* __cplusplus */

