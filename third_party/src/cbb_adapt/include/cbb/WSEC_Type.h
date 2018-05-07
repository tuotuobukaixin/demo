/*******************************************************************************
* Copyright @ Huawei Technologies Co., Ltd. 1998-2014. All rights reserved.  
* File name:  WSEC_Type.h
* Decription: 
     类型定义
*********************************************************************************/
#ifndef __WIRELESS_TYPE_D13A02_DCRFASD3F_4E7HLPO4C27
#define __WIRELESS_TYPE_D13A02_DCRFASD3F_4E7HLPO4C27

#include <stdio.h>
#include <time.h>

#ifdef __cplusplus
extern "C"
{
#endif

/************************************************************************
 *  macro defines
************************************************************************/

#define INOUT

#define WSEC_NULL_PTR       ((void *)0)

#define WSEC_FALSE          (0)
#define WSEC_TRUE           (1)
#define WSEC_FALSE_FOREVER  (!__LINE__) /* 恒假, 用于条件判断避免编译告警 */
#define WSEC_TRUE_FOREVER   (__LINE__)  /* 恒真, 用于条件判断避免编译告警 */

/************************************************************************
 *  type defines
************************************************************************/

/* basic type redefines */
typedef void            WSEC_VOID;
typedef unsigned char   WSEC_UINT8;
typedef unsigned char   WSEC_BYTE;
typedef unsigned short  WSEC_UINT16;
typedef unsigned int    WSEC_UINT32;

typedef int             WSEC_INT32;
typedef int             WSEC_BOOL;
typedef char            WSEC_CHAR;

typedef void *          WSEC_HANDLE;
typedef FILE*           WSEC_FILE;
typedef int             WSEC_POSITION;
typedef unsigned long   WSEC_FILE_LEN;

typedef unsigned int    WSEC_SIZE_T;
typedef clock_t         WSEC_CLOCK_T;
typedef unsigned long   WSEC_ERR_T;

/* 通知APP的通告字(以此识别是哪种通告) */
typedef enum
{
    WSEC_KMC_NTF_RK_EXPIRE          = 1, /* Root Key(RK)密钥物料(即将)过期 */
    WSEC_KMC_NTF_MK_EXPIRE          = 2, /* Master Key(MK)(即将)过期 */
    WSEC_KMC_NTF_MK_CHANGED         = 3, /* MK变更 */
    WSEC_KMC_NTF_USING_EXPIRED_MK   = 4, /* 使用过期MK */
    WSEC_KMC_NTF_KEY_STORE_CORRUPT  = 5, /* Keystore被破坏 */
    WSEC_KMC_NTF_CFG_FILE_CORRUPT   = 6, /* KMC配置文件被破坏 */
    WSEC_KMC_NTF_WRI_KEY_STORE_FAIL = 7, /* 保存Keystore失败 */
    WSEC_KMC_NTF_WRI_CFG_FILE_FAIL  = 8, /* 保存KMC配置文件失败 */
    WSEC_KMC_NTF_MK_NUM_OVERFLOW    = 9  /* MK个数即将超规格 */
} WSEC_NTF_CODE_ENUM;

/* CBB向APP通告 */
typedef WSEC_VOID (*WSEC_FP_Notify)(WSEC_NTF_CODE_ENUM eNtfCode, const WSEC_VOID* pData, WSEC_SIZE_T nDataSize);

/* 进度上报 */
typedef WSEC_VOID (*WSEC_FP_RptProgress)(WSEC_UINT32 ulTag, WSEC_UINT32 ulScale, WSEC_UINT32 ulCurrent, WSEC_BOOL* pbCancel);

/* struct defines */
#pragma pack(1)
typedef struct wsectagSysTime
{
    WSEC_UINT16 uwYear;     /* 年份 */
    WSEC_UINT8  ucMonth;    /* 月份(1~12) */
    WSEC_UINT8  ucDate;     /* 日期(1~31, 上限由年月份确定) */
    WSEC_UINT8  ucHour;     /* 时(0~23) */
    WSEC_UINT8  ucMinute;   /* 分(0~59) */
    WSEC_UINT8  ucSecond;   /* 秒(0~59) */
    WSEC_UINT8  ucWeek;     /* 星期(1~7分别表示周一到周日) */
} WSEC_SYSTIME_T;
#pragma pack()

#pragma pack(1)
typedef struct tagWSEC_SCHEDULE_TIME
{
    WSEC_UINT8  ucHour;   /* 时(0~23) */
    WSEC_UINT8  ucMinute; /* 分(0~59) */
    WSEC_UINT8  ucWeek;   /* 每周几, 取值为1~7, 表示周1~周日; 1~7之外，则表示每日 */
    WSEC_BYTE   abReserved[4]; /* 预留 */
} WSEC_SCHEDULE_TIME_STRU;
#pragma pack()

typedef struct tagWSEC_PROGRESS_RPT
{
    WSEC_UINT32         ulTag;         /* APP识别大任务的标签 */
    WSEC_FP_RptProgress pfRptProgress; /* 进度上报回调函数 */
} WSEC_PROGRESS_RPT_STRU; /* CBB处理大任务时向APP上报进度 */

#ifdef __cplusplus
}
#endif  /* __cplusplus */

#endif/* __WIRELESS_TYPE_D13A02_DCRFASD3F_4E7HLPO4C27 */
