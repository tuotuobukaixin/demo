/******************************************************************************

版权所有 (C), 2001-2011, 华为技术有限公司

******************************************************************************
文 件 名   : KMC_Itf.h
版 本 号   : 初稿
作    者   : 
生成日期   : 2014年6月16日
最近修改   :
功能描述   : KMC_Func.c 的对外接口头文件
函数列表   :
修改历史   :
1.日    期   : 2014年6月16日
作    者   : 
修改内容   : 创建文件

******************************************************************************/
#ifndef __KMC_ITF_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__
#define __KMC_ITF_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__

#include "WSEC_Type.h"
#include "WSEC_Config.h"

#ifdef __cplusplus
#if __cplusplus
extern "C"
{
#endif
#endif /* __cplusplus */

/*==============================================
                宏定义
==============================================*/
#define WSEC_MK_LEN_MAX        (128)  /* MK密钥明(密)文最大长度(字节数) */
#define KMC_KEY_DFT_ITERATIONS (4000) /* 缺省密钥迭代次数 */

/*==============================================
                枚举类型
==============================================*/
/* 以用途对密钥分类，按bit定义 */
typedef enum
{
    KMC_KEY_TYPE_ENCRPT    = 1,       /* 对称加密 */
    KMC_KEY_TYPE_INTEGRITY = 2,       /* 完整性保护 */
    KMC_KEY_TYPE_ENCRPT_INTEGRITY = 3 /* 加密及完整性 */
} KMC_KEY_TYPE_ENUM;

/* 数据保护模块提供的数据保护算法类型 */
typedef enum 
{
    SDP_ALG_ENCRPT,     /* 加密 */
    SDP_ALG_INTEGRITY,  /* 完整性保护 */
    SDP_ALG_PWD_PROTECT /* 口令保护 */
} KMC_SDP_ALG_TYPE_ENUM;

/* Root Key物料产生方式 */
typedef enum
{
    KMC_RK_GEN_BY_INNER, /* 系统自动生成 */
    KMC_RK_GEN_BY_IMPORT /* 外部导入 */
} KMC_RK_GEN_FROM;

/* Master Key产生方式 */
typedef enum
{
    KMC_MK_GEN_BY_INNER, /* 系统自动生成 */
    KMC_MK_GEN_BY_IMPORT /* 外部导入 */
} KMC_MK_GEN_FROM;

/* 密钥状态 */
typedef enum
{
    KMC_KEY_STATUS_INACTIVE = 0, /* 非活动状态的密钥不再用于机密数据加密，但可以用来解密历史密文 */
    KMC_KEY_STATUS_ACTIVE        /* 正常使用中 */
} KMC_KEY_STATUS_ENUM;

/* 密钥变更类型 */
typedef enum 
{
    KMC_KEY_ACTIVATED = 0, /* 密钥激活 */
    KMC_KEY_INACTIVATED,   /* 密钥去激活(过期) */
    KMC_KEY_REMOVED        /* 密钥被删除 */
} KMC_KEY_CHANGE_TYPE_ENUM;

/*==============================================
                结构体
==============================================*/
/*----------------------------------------------------------
1. Root Key(RK)信息
----------------------------------------------------------*/
typedef struct tagKMC_RK_ATTR
{
    WSEC_UINT16    usVer;              /* 版本号 */
    WSEC_UINT16    usRkMeterialFrom;   /* 根密钥物料来源, 见 KMC_RK_GEN_FROM */
    WSEC_SYSTIME_T stRkCreateTimeUtc;  /* 根密钥创建时间(UTC) */
    WSEC_SYSTIME_T stRkExpiredTimeUtc; /* 根密钥过期时间(UTC) */
    WSEC_UINT32    ulRmkIterations;    /* 派生RMK迭代次数 */
} KMC_RK_ATTR_STRU;

/*----------------------------------------------------------
2. Master Key(MK)信息
----------------------------------------------------------*/
#pragma pack(1)
/*``````````````````````````````````````````````````````````````````````````````````````````````
MK有两类关键字:
(1) ulDomainId + ulKeyId为唯一关键字, 用于识别MK
(2) ulDomainId + usType + ucStatus 为可重复关键字, 用于APP按usType获取当前状态为'可用'的MK
``````````````````````````````````````````````````````````````````````````````````````````````*/
typedef struct tagKMC_MK_INFO
{
    WSEC_UINT32    ulDomainId;      /* 密钥作用域 */
    WSEC_UINT32    ulKeyId;         /* 密钥ID，在同一Domain唯一 */
    WSEC_UINT16    usType;          /* 以用途对密钥分类, 见 KMC_KEY_TYPE_ENUM */
    WSEC_UINT8     ucStatus;        /* 密钥状态, 见 KMC_KEY_STATUS_ENUM */
    WSEC_UINT8     ucGenStyle;      /* 密钥产生方式, 见 KMC_MK_GEN_FROM */
    WSEC_SYSTIME_T stMkCreateTimeUtc;  /* MK创建时间(UTC) */
    WSEC_SYSTIME_T stMkExpiredTimeUtc; /* MK过期时间(UTC) */
} KMC_MK_INFO_STRU; /* MK头信息 */
#pragma pack()

/*----------------------------------------------------------
3. 密钥管理相关的配置
----------------------------------------------------------*/
/* 1) 全局性的密钥配置信息 */
#pragma pack(1)
typedef struct tagKMC_CFG_ROOT_KEY
{
    WSEC_UINT32    ulRootKeyLifeDays;          /* Rootkey 有效时间（天）*/
    WSEC_UINT32    ulRootMasterKeyIterations;  /* Rootkey 迭代次数 */
    WSEC_BYTE      abReserved[8];              /* 预留 */
} KMC_CFG_ROOT_KEY_STRU; /* RK管理参数 */
#pragma pack()
#pragma pack(1)
typedef struct tagKMC_CFG_KEY_MAN
{
    WSEC_UINT32    ulWarningBeforeKeyExpiredDays; /* 密钥过期提前预警天数 */
    WSEC_UINT32    ulGraceDaysForUseExpiredKey;   /* 使用过期密钥的宽限天数, 超过宽限期则通知APP */
    WSEC_BOOL      bKeyAutoUpdate;                /* 密钥过期是否自动更新, 针对根密钥、及KeyFrom=0的MK有效 */
    WSEC_SCHEDULE_TIME_STRU stAutoUpdateKeyTime;  /* 密钥自动更新时间 */
    WSEC_BYTE      abReserved[8];                 /* 预留 */
} KMC_CFG_KEY_MAN_STRU; /* 密钥管理配置参数 */
#pragma pack()

/* 2) 数据保护 */
#pragma pack(1)
typedef struct tagKMC_CFG_DATA_PROTECT
{
    WSEC_UINT32 ulAlgId;         /* 算法ID */
    WSEC_UINT16 usKeyType;       /* 以用途对密钥分类, 见 KMC_KEY_TYPE_ENUM */
    WSEC_BOOL   bAppendMac;      /* 是否追加完整性校验值 */
    WSEC_UINT32 ulKeyIterations; /* 密钥迭代次数 */
    WSEC_BYTE   abReserved[8];   /* 预留 */
} KMC_CFG_DATA_PROTECT_STRU;
#pragma pack()

/* 3) DOMAIN Key Type配置 */
#pragma pack(1)
typedef struct tagKMC_CFG_KEY_TYPE
{
    WSEC_UINT16 usKeyType;     /* 以用途对密钥分类, 见 KMC_KEY_TYPE_ENUM */
    WSEC_UINT32 ulKeyLen;      /* 密钥长度 */
    WSEC_UINT32 ulKeyLifeDays; /* 有效期(天数) */
    WSEC_BYTE   abReserved[8]; /* 预留 */
} KMC_CFG_KEY_TYPE_STRU;
#pragma pack()

/* 4) DOMAIN配置 */
#pragma pack(1)
typedef struct tagKMC_CFG_DOMAIN_INFO
{
    WSEC_UINT32  ulId;        /* 密钥作用域 */
    WSEC_UINT8   ucKeyFrom;   /* 密钥生成来源, 见 KMC_MK_GEN_FROM */
    WSEC_CHAR    szDesc[128]; /* 密钥描述 */
    WSEC_BYTE    abReserved[8]; /* 预留 */
} KMC_CFG_DOMAIN_INFO_STRU;
#pragma pack()

/*----------------------------------------------------------
4. KMC所需的文件名
----------------------------------------------------------*/
typedef struct tagKMC_FILE_NAME
{
    WSEC_CHAR* pszKeyStoreFile[2]; /* Keystore文件名(可靠性考虑, 两文件互为主备) */
    WSEC_CHAR* pszKmcCfgFile[2];   /* KMC配置文件名(如果APP自己管理配置数据, 则不需提供) */
} KMC_FILE_NAME_STRU;

/*----------------------------------------------------------
5. 通知APP的数据结构
----------------------------------------------------------*/
/* 1) RK即将过期通告 */
typedef struct tagKMC_RK_EXPIRE_NTF
{
    KMC_RK_ATTR_STRU stRkInfo;    /* 即将过期的根密钥信息 */
    WSEC_INT32       nRemainDays; /* 离过期日还有多少天 */
} KMC_RK_EXPIRE_NTF_STRU;

/* 2) MK即将过期通告 */
typedef struct tagKMC_MK_EXPIRE_NTF
{
    KMC_MK_INFO_STRU stMkInfo;    /* 即将过期的MK信息 */
    WSEC_INT32       nRemainDays; /* 离过期日还有多少天 */
} KMC_MK_EXPIRE_NTF_STRU;

/* 3) MK变更通告 */
typedef struct tagKMC_MK_CHANGE_NTF
{
    KMC_MK_INFO_STRU         stMkInfo; /* 变更的MK信息 */
    KMC_KEY_CHANGE_TYPE_ENUM eType;    /* 变更类型 */
} KMC_MK_CHANGE_NTF_STRU;

/* 4) 过期MK超宽限期使用通告 */
typedef struct tagKMC_USE_EXPIRED_MK_NTF
{
    KMC_MK_INFO_STRU stExpiredMkInfo; /* 过期MK信息 */
    WSEC_INT32       nExpiredDays;    /* 过期天数 */
} KMC_USE_EXPIRED_MK_NTF_STRU;

/* 5) 写Keystore文件失败通告 */
typedef struct tagKMC_WRI_KSF_FAIL_NTF
{
    WSEC_ERR_T ulCause; /* 失败原因 */
} KMC_WRI_KSF_FAIL_NTF_STRU;

/* 6) 写KMC配置文件失败通告 */
typedef struct tagKMC_WRI_KCF_FAIL_NTF
{
    WSEC_ERR_T ulCause; /* 失败原因 */
} KMC_WRI_KCF_FAIL_NTF_STRU;

/* 7) MK个数超规格通告 */
typedef struct tagKMC_MK_NUM_OVERFLOW
{
    WSEC_UINT32 ulNum;    /* 当前MK个数 */
    WSEC_UINT32 ulMaxNum; /* 允许MK最大个数 */
} KMC_MK_NUM_OVERFLOW_STRU;

/*----------------------------------------------
    读配置回调函数指针定义
----------------------------------------------*/
typedef WSEC_BOOL (*WSEC_FP_ReadRootKeyCfg)(KMC_CFG_ROOT_KEY_STRU* pstRkCfg); /* 读取RootKey配置 */
typedef WSEC_BOOL (*WSEC_FP_ReadKeyManCfg)(KMC_CFG_KEY_MAN_STRU* pstKmCfg); /* 读取KEY管理配置 */

/* 读取所有Domain个数 */
typedef WSEC_BOOL (*WSEC_FP_ReadCfgOfDomainCount)(WSEC_UINT32* pulDomainCount); /* 读取KMC配置数据之 Domain个数 */

/* 读取所有Domain配置信息 */
typedef WSEC_BOOL (*WSEC_FP_ReadCfgOfDomainInfo)(KMC_CFG_DOMAIN_INFO_STRU* pstAllDomainInfo, /* 用于输出所有Domain配置信息的的缓冲区 */
                                                 WSEC_UINT32 ulDomainCount); /* pstAllDomainInfo所指向缓冲区容纳KMC_CFG_DOMAIN_INFO_STRU数据结构的个数 */

/* 读取指定Domain有多少条KeyType配置 */
typedef WSEC_BOOL (*WSEC_FP_ReadCfgOfDomainKeyTypeCount)(WSEC_UINT32 ulDomainId, /* 指定的Domain */
                                                       WSEC_UINT32* pulKeyTypeCount); /* 输出该Domain的KeyType记录数 */

/* 读取指定Domain的所有KeyType配置 */
typedef WSEC_BOOL (*WSEC_FP_ReadCfgOfDomainKeyType)(WSEC_UINT32 ulDomainId, /* 给定Doamin */
                                                    KMC_CFG_KEY_TYPE_STRU* pstDomainAllKeyType, /* 调用者给定的用于输出该Domain所有KeyType配置记录 */
                                                    WSEC_UINT32 ulKeyTypeCount); /* pstDomainAllKeyType所指向缓冲区容纳KMC_KEY_TYPE_STRU数据结构的个数 */

/* 读取指定类型的算法配置 */
typedef WSEC_BOOL (*WSEC_FP_ReadCfgOfDataProtection)(KMC_SDP_ALG_TYPE_ENUM eType, KMC_CFG_DATA_PROTECT_STRU *pstPara); 

/* KMC配置回调函数集 */
typedef struct tagKMC_FP_CFG_CALLBACK
{
    WSEC_FP_ReadRootKeyCfg              pfReadRootKeyCfg;              /* 读取RootKey配置相关参数 */
    WSEC_FP_ReadKeyManCfg               pfReadKeyManCfg;               /* 读取密钥生命周期管理相关参数 */
    WSEC_FP_ReadCfgOfDomainCount        pfReadCfgOfDomainCount;        /* 读取KMC配置数据之 Domain个数 */
    WSEC_FP_ReadCfgOfDomainInfo         pfReadCfgOfDomainInfo;         /* 读取所有Domain配置信息 */
    WSEC_FP_ReadCfgOfDomainKeyTypeCount pfReadCfgOfDomainKeyTypeCount; /* 读取指定Domain有多少条KeyType配置 */
    WSEC_FP_ReadCfgOfDomainKeyType      pfReadCfgOfDomainKeyType;      /* 读取指定Domain的所有KeyType配置 */
    WSEC_FP_ReadCfgOfDataProtection     pfReadCfgOfDataProtection;     /* 读取指定类型的算法配置 */
} KMC_FP_CFG_CALLBACK_STRU;

typedef struct tagKMC_FP_CALLBACK
{
    KMC_FP_CFG_CALLBACK_STRU          stReadCfg;     /* KMC配置回调函数集 */
} KMC_FP_CALLBACK_STRU; /* KMC的回调函数 */

/*----------------------------------------------
    函数原型说明
----------------------------------------------*/
/* (1) Keystore访问 */
WSEC_ERR_T KMC_RmvMk(WSEC_UINT32 ulDomainId, WSEC_UINT32 ulKeyId);
WSEC_ERR_T KMC_ExportMkFile(const WSEC_CHAR* pszToFile, const WSEC_BYTE* pbPwd, WSEC_UINT32 ulPwdLen, WSEC_UINT32 ulKeyIterations, const WSEC_PROGRESS_RPT_STRU* pstRptProgress);
WSEC_ERR_T KMC_ImportMkFile(const WSEC_CHAR* pszFromFile, const WSEC_BYTE* pbPwd, WSEC_UINT32 ulPwdLen, const WSEC_PROGRESS_RPT_STRU* pstRptProgress);
WSEC_ERR_T KMC_UpdateRootKey(const WSEC_BYTE* pbKeyEntropy, WSEC_SIZE_T ulSize);
WSEC_ERR_T KMC_GetRootKeyInfo(KMC_RK_ATTR_STRU* pstRkInfo);
WSEC_INT32 KMC_GetMkCount();
WSEC_ERR_T KMC_GetMk(WSEC_INT32 Index, KMC_MK_INFO_STRU* pstMk);
WSEC_ERR_T KMC_GetMaxMkId(WSEC_UINT32 ulDomainId, WSEC_UINT32* pulMaxKeyId);
WSEC_ERR_T KMC_SetMkExpireTime(WSEC_UINT32 ulDomainId, WSEC_UINT32 ulKeyId, const WSEC_SYSTIME_T* psExpireTime);
WSEC_ERR_T KMC_SetMkStatus(WSEC_UINT32 ulDomainId, WSEC_UINT32 ulKeyId, WSEC_UINT8 ucStatus);
WSEC_ERR_T KMC_RegisterMk(WSEC_UINT32 ulDomainId, WSEC_UINT32 ulKeyId, WSEC_UINT16 usKeyType, const WSEC_BYTE* pPlainTextKey, WSEC_UINT32 ulKeyLen);
WSEC_ERR_T KMC_CreateMk(WSEC_UINT32 ulDomainId, WSEC_UINT16 usKeyType);
WSEC_ERR_T KMC_GetMkDetail(WSEC_UINT32 ulDomainId, WSEC_UINT32 ulKeyId, KMC_MK_INFO_STRU* pstMkInfo, WSEC_BYTE* pbKeyPlainText, WSEC_UINT32* pKeyLen);
WSEC_ERR_T KMC_SecureEraseKeystore();

/* (2) 配置 */
/* 2.1 密钥属性配置 */
WSEC_ERR_T KMC_SetRootKeyCfg(const KMC_CFG_ROOT_KEY_STRU* pstRkCfg);
WSEC_ERR_T KMC_SetKeyManCfg(const KMC_CFG_KEY_MAN_STRU* pstKmCfg);
WSEC_ERR_T KMC_GetRootKeyCfg(KMC_CFG_ROOT_KEY_STRU* pstRkCfg);
WSEC_ERR_T KMC_GetKeyManCfg(KMC_CFG_KEY_MAN_STRU* pstKmCfg);

/* 2.2 Domain配置 */
WSEC_ERR_T KMC_AddDomain(const KMC_CFG_DOMAIN_INFO_STRU* pstDomain);
WSEC_ERR_T KMC_RmvDomain(WSEC_UINT32 ulDomainId);

WSEC_ERR_T KMC_AddDomainKeyType(WSEC_UINT32 ulDomainId, const KMC_CFG_KEY_TYPE_STRU* pstKeyType);
WSEC_ERR_T KMC_RmvDomainKeyType(WSEC_UINT32 ulDomainId, WSEC_UINT16 usKeyType);

WSEC_INT32 KMC_GetDomainCount(); /* 获取配置的Domain个数 */
WSEC_ERR_T KMC_GetDomain(WSEC_INT32 Index, KMC_CFG_DOMAIN_INFO_STRU* pstDomainInfo); /* 获取指定位置上的Domain */

WSEC_INT32 KMC_GetDomainKeyTypeCount(WSEC_UINT32 ulDomainId); /* 获取Domain下KeyType个数 */
WSEC_ERR_T KMC_GetDomainKeyType(WSEC_UINT32 ulDomainId, WSEC_INT32 Index, KMC_CFG_KEY_TYPE_STRU* pstKeyType); /* 获取指定Domain指定位置上的KeyType */

WSEC_VOID KMC_GetExpiredMkStartPos(WSEC_POSITION* pPos); /* 获取过期MK的首位置 */
WSEC_BOOL KMC_GetExpiredMkByPos(WSEC_POSITION* pPosNow, KMC_MK_INFO_STRU* pstExpiredMk); /* 获取指定位置上的过期MK并返回下一个过期MK的位置 */

WSEC_ERR_T KMC_SetDataProtectCfg(KMC_SDP_ALG_TYPE_ENUM eType, const KMC_CFG_DATA_PROTECT_STRU* pstPara);
WSEC_ERR_T KMC_GetDataProtectCfg(KMC_SDP_ALG_TYPE_ENUM eType, KMC_CFG_DATA_PROTECT_STRU *pstPara);

/* 列出CBB支持的安全算法 */
typedef WSEC_VOID (*KMC_FP_ProcAlg)(WSEC_UINT32 ulAlgId, const WSEC_CHAR* pszAlgName, WSEC_VOID* pReserved);
WSEC_ERR_T KMC_GetAlgList(KMC_FP_ProcAlg pfProcAlg, INOUT WSEC_VOID *pReserved);

#ifdef __cplusplus
#if __cplusplus
}
#endif
#endif /* __cplusplus */

#endif /* __KMC_ITF_H_D13DA0FG2_DCRFKLAPSD32SF_4EHLPOC27__ */

