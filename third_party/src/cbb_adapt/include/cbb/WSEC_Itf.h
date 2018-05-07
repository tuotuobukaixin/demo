/*******************************************************************************
* Copyright @ Huawei Technologies Co., Ltd. 1998-2014. All rights reserved.  
* File name: WSEC_Itf.h
* Decription:
  对外接口
*********************************************************************************/
#ifndef __WIRELESS_SEC_ITF_D13ASA042_DCRF3F_4E74C27
#define __WIRELESS_SEC_ITF_D13ASA042_DCRF3F_4E74C27

#include <stdio.h>
#include "WSEC_Config.h"
#include "WSEC_Type.h"
#include "KMC_Itf.h"
#include "WSEC_ErrorCode.h"

#ifdef __cplusplus
extern "C" {
#endif

/*算法ID*/

#define WSEC_ALGID_NUM_PER_TYPE    (1024)

#define WSEC_ALGID_SYM_BEGIN       (1 + WSEC_ALGID_NUM_PER_TYPE * 0)
#define WSEC_ALGID_DIGEST_BEGIN    (1 + WSEC_ALGID_NUM_PER_TYPE * 1)
#define WSEC_ALGID_HMAC_BEGIN      (1 + WSEC_ALGID_NUM_PER_TYPE * 2)
#define WSEC_ALGID_PBKDF_BEGIN     (1 + WSEC_ALGID_NUM_PER_TYPE * 3)

typedef enum
{
    WSEC_ALGID_UNKNOWN,         /* Unknown alg id */

    /*************** Symmetric Algorithm ****************/

    WSEC_ALGID_DES_EDE3_ECB =  WSEC_ALGID_SYM_BEGIN + 1,   /* identifies 3 key triple DES algorithm  ECB mode */
    WSEC_ALGID_DES_EDE3_CBC,   /* identifies 3 key triple DES algorithm  CBC mode */

    WSEC_ALGID_AES128_ECB,     /* identifies AES-128 algorithm  ECB mode */
    WSEC_ALGID_AES128_CBC,     /* identifies AES-128 algorithm  CBC mode */

    WSEC_ALGID_AES256_ECB,     /* identifies AES-256 algorithm  ECB mode */
    WSEC_ALGID_AES256_CBC,     /* identifies AES-256 algorithm  CBC mode */

    /*************** Digest Algorithm ****************/
    WSEC_ALGID_MD5 =   WSEC_ALGID_DIGEST_BEGIN, 
                       /* identifies the MD5 hash algorithm */
    WSEC_ALGID_SHA1,   /* identifies the SHA1 hash algorithm */
    WSEC_ALGID_SHA224, /* identifies the SHA224 hash algorithm */
    WSEC_ALGID_SHA256, /* identifies the SHA256 hash algorithm */
    WSEC_ALGID_SHA384, /* identifies the SHA384 hash algorithm */
    WSEC_ALGID_SHA512, /* identifies the SHA512 hash algorithm */

    /*************** HMAC Algorithm ****************/
    WSEC_ALGID_HMAC_MD5 =      WSEC_ALGID_HMAC_BEGIN,
                               /* identifies hmac with MD5 */
    WSEC_ALGID_HMAC_SHA1,      /* identifies hmac with SHA1 */
    WSEC_ALGID_HMAC_SHA224,    /* identifies hmac with SHA224 */
    WSEC_ALGID_HMAC_SHA256,    /* identifies hmac with SHA256 */
    WSEC_ALGID_HMAC_SHA384,    /* identifies hmac with SHA384 */
    WSEC_ALGID_HMAC_SHA512,    /* identifies hmac with SHA512 */

    /*************** PBKDF Algorithm ****************/
    WSEC_ALGID_PBKDF2_HMAC_MD5 =      WSEC_ALGID_PBKDF_BEGIN,
                                      /* identifies hmac with MD5 used in pbkdf2 */
    WSEC_ALGID_PBKDF2_HMAC_SHA1,      /* identifies hmac with SHA1 used in pbkdf2 */
    WSEC_ALGID_PBKDF2_HMAC_SHA224,    /* identifies hmac with SHA224 used in pbkdf2 */
    WSEC_ALGID_PBKDF2_HMAC_SHA256,    /* identifies hmac with SHA256 used in pbkdf2 */
    WSEC_ALGID_PBKDF2_HMAC_SHA384,    /* identifies hmac with SHA384 used in pbkdf2 */
    WSEC_ALGID_PBKDF2_HMAC_SHA512,    /* identifies hmac with SHA512 used in pbkdf2 */

    //WSEC_ALGID_PKCS8,
}WSEC_ALGID_E;


#define WSEC_FILEPATH_MAX_LEN (260) /* 文件路径最大长度 */

/************************************************************************
 *  function pointer                                                    
************************************************************************/
/* write log */
typedef WSEC_VOID  (*WSEC_FP_WriLog)(int nLevel, const char* pszModuleName, const char* pszOccurFileName, int nOccurLine, const char* pszLog);

/* memory operation */
typedef WSEC_VOID* (*WSEC_FP_MemAlloc)( WSEC_SIZE_T uSize);
typedef WSEC_VOID  (*WSEC_FP_MemFree)( WSEC_VOID* pvMem);
typedef WSEC_INT32 (*WSEC_FP_MemCmp)(const WSEC_VOID *Buf1, const WSEC_VOID *Buf2, WSEC_SIZE_T Count);

/* Lock */
typedef WSEC_BOOL (*WSEC_FP_CreateLock)( WSEC_HANDLE *phMutex);
typedef WSEC_VOID (*WSEC_FP_DestroyLock)( WSEC_HANDLE  hMutex);
typedef WSEC_VOID (*WSEC_FP_Lock)( WSEC_HANDLE  hMutex);
typedef WSEC_VOID (*WSEC_FP_Unlock)( WSEC_HANDLE  hMutex);

typedef WSEC_VOID (*WSEC_FP_DoEvents)();

/* file IO operation */
typedef WSEC_FILE (*WSEC_FP_Fopen)(const char *filename, const char *mode);
typedef int (*WSEC_FP_Fclose)(WSEC_FILE stream);
typedef size_t (*WSEC_FP_Fread)(void *buffer, size_t size, size_t count, WSEC_FILE stream);
typedef size_t (*WSEC_FP_Fwrite)(const void *buffer, size_t size, size_t count, WSEC_FILE stream);
typedef int (*WSEC_FP_Fflush)(WSEC_FILE stream);
typedef int (*WSEC_FP_Fremove)(const char *path);
typedef int (*WSEC_FP_Fgetc)(WSEC_FILE stream);
typedef char* (*WSEC_FP_Fgets)(char *string, int n, WSEC_FILE stream);
typedef long (*WSEC_FP_Ftell)(WSEC_FILE stream);
typedef int (*WSEC_FP_Fseek)(WSEC_FILE stream, long offset, int origin);
typedef int (*WSEC_FP_Feof)(WSEC_FILE stream);
typedef int (*WSEC_FP_Ferror)(WSEC_FILE stream);

/* 文件时间设置/获取回调函数 */
typedef WSEC_BOOL (*WSEC_FP_GetFileDateTime)(const WSEC_CHAR* pszFileName, WSEC_SYSTIME_T* pstCreateTime, WSEC_SYSTIME_T* pstLastEditTime);
typedef WSEC_BOOL (*WSEC_FP_SetFileDateTime)(const WSEC_CHAR* pszFileName, const WSEC_SYSTIME_T* pstCreateTime, const WSEC_SYSTIME_T* pstLastEditTime);

/* 向APP输出CBB结构长度, 调测用 */
typedef WSEC_VOID (*WSEC_FP_ShowStructSize)(const WSEC_CHAR* pszStructName, WSEC_SIZE_T nSize);

/************************************************************************
 *  struct define                                                    
************************************************************************/
/* 系统基本回调函数结构 */
typedef struct tagWSEC_MEMORY_CALLBACK_STRU
{
    WSEC_FP_MemAlloc            pfMemAlloc;
    WSEC_FP_MemFree             pfMemFree;
    WSEC_FP_MemCmp              pfMemCmp;
} WSEC_MEMORY_CALLBACK_STRU;

/* 系统文件读写回调函数结构 */
typedef struct tagWSEC_FILE_CALLBACK_STRU
{
    WSEC_FP_Fopen               pfFopen;
    WSEC_FP_Fclose              pfFclose;
    WSEC_FP_Fread               pfFread;
    WSEC_FP_Fwrite              pfFwrite;
    WSEC_FP_Fflush              pfFflush;
    WSEC_FP_Fremove             pfFremove;        
    WSEC_FP_Fgetc               pfFgetc;
    WSEC_FP_Fgets               pfFgets;
    WSEC_FP_Ftell               pfFtell;
    WSEC_FP_Fseek               pfFseek;
    WSEC_FP_Feof                pfFeof;
    WSEC_FP_Ferror              pfFerror;
} WSEC_FILE_CALLBACK_STRU;

/* 系统线程回调函数结构 */
typedef struct tagWSEC_LOCK_CALLBACK_STRU
{
    WSEC_FP_CreateLock          pfCreateLock;
    WSEC_FP_DestroyLock         pfDestroyLock;
    WSEC_FP_Lock                pfLock;
    WSEC_FP_Unlock              pfUnlock;
}WSEC_LOCK_CALLBACK_STRU;

typedef struct tagWSEC_BASE_RELY_APP_CALLBACK_STRU
{
    WSEC_FP_WriLog              pfWriLog;
    WSEC_FP_Notify              pfNotify;
    WSEC_FP_DoEvents            pfDoEvents;
} WSEC_BASE_RELY_APP_CALLBACK_STRU;

typedef struct tagWSEC_FP_CALLBACK
{
    WSEC_MEMORY_CALLBACK_STRU        stMemory;
    WSEC_FILE_CALLBACK_STRU          stFile;
    WSEC_LOCK_CALLBACK_STRU          stLock;
    WSEC_BASE_RELY_APP_CALLBACK_STRU stRelyApp;
    KMC_FP_CALLBACK_STRU             stKmcCallbackFun;
} WSEC_FP_CALLBACK_STRU;

/************************************************************************
 *  interface define                                                    
************************************************************************/
/* 系统启动或关闭时，APP需要分别调用如下函数 */
WSEC_ERR_T WSEC_Initialize(const KMC_FILE_NAME_STRU* pstFileName, /* CBB需要的文件名, 若使用SDP或KMC则必须提供，否则可NULL */
                           const WSEC_FP_CALLBACK_STRU* pstCallbackFun, /* CBB回调函数 */
                           const WSEC_PROGRESS_RPT_STRU* pstRptProgress, /* 初始化过程中如果耗时较长则上报进度 */
                           const WSEC_VOID* pvReserved); /* 保留参数 */
WSEC_VOID WSEC_OnTimer(const WSEC_PROGRESS_RPT_STRU* pstRptProgress);
WSEC_ERR_T WSEC_Finalize();
WSEC_ERR_T WSEC_Reset();
const WSEC_CHAR* WSEC_GetVerion();

#ifdef __cplusplus
}
#endif  /* __cplusplus */

#endif/* __WIRELESS_SEC_ITF_D13ASA042_DCRF3F_4E74C27 */
