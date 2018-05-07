/******************************************************************************

                  版权所有 (C), 2001-2011, 华为技术有限公司

 ******************************************************************************
  文 件 名   : SDP_Pri.h
  版 本 号   : 初稿
  作    者   : 
  生成日期   : 2014年6月16日
  最近修改   :
  功能描述   : SDP_Func.c 的内部接口头文件，不对外开放
  函数列表   :
  修改历史   :
  1.日    期   : 2014年6月16日
    作    者   : 
    修改内容   : 创建文件

******************************************************************************/
#ifndef __SDP_PRI_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__
#define __SDP_PRI_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__

#include "WSEC_Type.h"
#include "WSEC_Pri.h"
#include "CAC_Pri.h"
#include "KMC_Itf.h"
#include "SDP_Itf.h"

#ifdef __cplusplus
#if __cplusplus
extern "C"{
#endif
#endif /* __cplusplus */

#define WSEC_HMAC_LEN_MAX (64)

/*----------------------------------------------*
 * 边界值定义
 *----------------------------------------------*/
#define SDP_SALT_LEN            16u
#define SDP_IV_MAX_LEN          16u /* 对称算法使用的IV的长度 */
#define SDP_PTMAC_MAX_LEN       64u /* 对称算法密文头包含明文HMAC的长度 */
#define SDP_KEY_MAX_LEN        128u
#define SDP_ALGNAME_MAX_LEN     64u /* 算法名称的长度 */

/*----------------------------------------------*
 * 密文文件TLV格式之Tag定义
 *----------------------------------------------*/
typedef enum 
{
    SDP_CFT_FILE_HDR = 1, /* 密文文件头 */
    SDP_CFT_CIPHER_HDR,   /* 密文头 */
    SDP_CFT_CIPHER_BODY,  /* 密文体 */
    SDP_CFT_HMAC_VAL      /* HMAC值 */
} SDP_CIPHER_FILE_TLV_ENUM;

/*----------------------------------------------*
 * 密文头定义
 *----------------------------------------------*/
#pragma pack(1)
typedef struct
{
    WSEC_UINT32 ulVersion; /* 数据保护模块版本号 */
    WSEC_BOOL   bHmacFlag; /* 是否包含明文HMAC,若包含明文HMAC则明文HMAC存放在密文的后面. 对于大数据的加解密暂时不支持明文HMAC*/
    WSEC_UINT32 ulDomain;  /* KEYID对应的域 */
    WSEC_UINT32 ulAlgId;   /* 算法ID*/
    WSEC_UINT32 ulKeyId;/* 计算HMAC使用的KEYID */
    WSEC_UINT32 ulIterCount; /* 迭代轮次，使用应用配置给密钥管理模块派生工作密钥的轮次 */
    WSEC_UINT8 aucSalt[SDP_SALT_LEN];/* 数据保护模块生成的盐值 */
    WSEC_UINT8 aucIV[SDP_IV_MAX_LEN];/* 数据保护模块生成的盐值 */    
    WSEC_UINT32 ulCDLen;/* 加密后的密文数据长度 */
}SDP_CIPHER_HEAD_STRU;/* 对称加密的数据头 */
#pragma pack()

/*----------------------------------------------*
 * HMAC头定义
 *----------------------------------------------*/
#pragma pack(1)
typedef struct
{
    WSEC_UINT32 ulVersion;  /* 数据保护模块版本号 */
    WSEC_UINT32 ulDomain;  /* KEYID对应的域 */
    WSEC_UINT32 ulAlgId;   /* 算法ID*/
    WSEC_UINT32 ulKeyId;/* 计算HMAC使用的KEYID */
    WSEC_UINT32 ulIterCount; /* 迭代轮次，使用应用配置给密钥管理模块派生工作密钥的轮次 */
    WSEC_UINT8 aucSalt[SDP_SALT_LEN];/* 数据保护模块生成的盐值 */
}SDP_HMAC_HEAD_STRU;
#pragma pack()

/*----------------------------------------------*
 * 口令加密头定义
 *----------------------------------------------*/
#pragma pack(1)
typedef struct
{
    WSEC_UINT32 ulVersion;  /* 数据保护模块版本号 */
    WSEC_UINT32 ulAlgId;    /* 算法ID*/
    WSEC_UINT32 ulIterCount; /* 迭代轮次，使用应用配置给密钥管理模块派生工作密钥的轮次 */
    WSEC_UINT8 aucSalt[SDP_SALT_LEN];/* 数据保护模块生成的盐值 */
    WSEC_UINT32 ulCDLen;     /* 加密后的口令长度 */
}SDP_PWD_HEAD_STRU;
#pragma pack()

/*----------------------------------------------*
 * 大数据加解密运算的上下文定义
 *----------------------------------------------*/
typedef struct
{
    WSEC_CRYPT_CTX          stWsecCtx;    /* CAC 适配层上下文环境 */
    SDP_CRYPT_CTX           stSdpCtxHmac; /* SDP 适配层上下文环境，用于完整性校验 */
    SDP_CIPHER_HEAD_STRU    stCipherHead; /* 大数据加解密的已记录密文头部 */
    SDP_HMAC_HEAD_STRU      stHmacHead;   /* 大数据计算HMAC的已记录HMAC头部 */
    WSEC_ALGTYPE_E          eAlgType;     /* 当前执行的算法类型 */
} SDP_CRYPT_CTX_STRU;

/*----------------------------------------------*
 * 错误上下文信息
 *----------------------------------------------*/
typedef struct
{
    WSEC_UINT32          ulErrCount;         /* 错误个数 */
    WSEC_UINT32          ulLastErrCode;      /* 上次的错误码 */
} SDP_ERROR_CTX_STRU;

/*----------------------------------------------*
 * 密文文件头结构定义
 *----------------------------------------------*/
#pragma pack(1)
typedef struct tagSDP_CIPHER_FILE_HDR
{
    WSEC_BYTE      abFormatFlag[32];    /* 格式标记符 */
    WSEC_UINT32    ulVer;               /* 密文文件版本号 */
    WSEC_UINT32    ulPlainBlockLenMax;  /* 最大明文段长度 */
    WSEC_UINT32    ulCipherBlockLenMax; /* 最大密文段长度 */
    WSEC_SYSTIME_T tCreateFileTimeUtc;  /* 密文文件生成时间(UTC) */
    WSEC_SYSTIME_T tSrcFileCreateTime;  /* 源文件创建时间 */
    WSEC_SYSTIME_T tSrcFileEditTime;    /* 源文件最近修改时间 */
    WSEC_BYTE      abReserved[16];      /* 预留 */
} SDP_CIPHER_FILE_HDR_STRU;
#pragma pack()

#define SDP_CIPHER_HEAD_STRU_LEN     sizeof(SDP_CIPHER_HEAD_STRU)
#define SDP_HMAC_HEAD_STRU_LEN       sizeof(SDP_HMAC_HEAD_STRU)
#define SDP_PWD_HEAD_STRU_LEN        sizeof(SDP_PWD_HEAD_STRU)
#define SDP_CRYPT_CTX_LEN            sizeof(SDP_CRYPT_CTX_STRU)

/*----------------------------------------------*
 * 检查长度是否超过保留长度限制                 *
 *----------------------------------------------*/
typedef WSEC_BYTE AssertValidLengthCipherHead[SDP_CIPHER_HEAD_LEN - SDP_CIPHER_HEAD_STRU_LEN];
typedef WSEC_BYTE AssertValidLengthHmacHead[SDP_HMAC_HEAD_LEN - SDP_HMAC_HEAD_STRU_LEN];
typedef WSEC_BYTE AssertValidLengthPwdHead[SDP_PWD_HEAD_LEN - SDP_PWD_HEAD_STRU_LEN];

/*----------------------------------------------*
 * 私有函数原型说明                             *
 *----------------------------------------------*/
WSEC_ERR_T SDP_GetAlgProperty(
    WSEC_UINT32 ulAlgID, WSEC_CHAR *pcAlgName, WSEC_UINT32 ulAlgNameLen,
    WSEC_ALGTYPE_E *peAlgType,
    WSEC_UINT32 *pulKeyLen,
    WSEC_UINT32 *pulIVLen,
    WSEC_UINT32 *pulBlockLen,
    WSEC_UINT32 *pulMACLen);

WSEC_VOID SDP_CvtByteOrder4CipherFileHdr(SDP_CIPHER_FILE_HDR_STRU* pstFileHdr, WSEC_BYTEORDER_CVT_ENUM eOper);
WSEC_VOID SDP_CvtByteOrder4CipherTextHeader(SDP_CIPHER_HEAD_STRU* pstHdr, WSEC_BYTEORDER_CVT_ENUM eOper);
WSEC_VOID SDP_CvtByteOrder4HmacTextHeader(SDP_HMAC_HEAD_STRU* pstHdr, WSEC_BYTEORDER_CVT_ENUM eOper);
WSEC_VOID SDP_CvtByteOrder4PwdCipherTextHeader(SDP_PWD_HEAD_STRU* pstHdr, WSEC_BYTEORDER_CVT_ENUM eOper);

WSEC_ERR_T SDP_FillCipherTextHeader(
    KMC_SDP_ALG_TYPE_ENUM eIntfType, WSEC_UINT32 ulDomain,
    SDP_CIPHER_HEAD_STRU *pstCipherHead,
    WSEC_BYTE *pucKey, WSEC_UINT32 *pulKeyLen,
    WSEC_UINT32 *pulIVLen);
WSEC_ERR_T SDP_FillHmacTextHeader(KMC_SDP_ALG_TYPE_ENUM eIntfType, WSEC_UINT32 ulDomain,
                                  SDP_HMAC_HEAD_STRU *pstHmacHead,
                                  WSEC_BYTE *pucKey, WSEC_UINT32 *pulKeyLen);
WSEC_ERR_T SDP_FillPwdCipherTextHeader(KMC_SDP_ALG_TYPE_ENUM eIntfType, SDP_PWD_HEAD_STRU *pstCipherHead);

WSEC_ERR_T SDP_GetWorkKey(
    WSEC_UINT32 ulDomain,
    WSEC_UINT16 usKeyType,
    WSEC_UINT32 *pulKeyId,
    WSEC_UINT32 *pulIterCount,
    WSEC_BYTE *pucSalt, WSEC_UINT32 ulSaltLen,
    WSEC_BYTE *pucIV, WSEC_UINT32 ulIVLen,
    WSEC_BYTE *pucKey, WSEC_UINT32 ulKeyLen);
WSEC_ERR_T SDP_GetWorkKeyByID(
    WSEC_UINT32 ulDomain,
    WSEC_UINT32 ulKeyId,
    WSEC_UINT32 ulIterCount,
    const WSEC_BYTE *pucSalt, WSEC_UINT32 ulSaltLen,
    WSEC_BYTE *pucKey, WSEC_UINT32 ulKeyLen);
WSEC_VOID SDP_FreeCtx(SDP_CRYPT_CTX *pstSdpCtx);

#ifdef __cplusplus
#if __cplusplus
}
#endif
#endif /* __cplusplus */


#endif /* __SDP_PRI_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__ */

