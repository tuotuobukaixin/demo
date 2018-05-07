/*******************************************************************************
* Copyright @ Huawei Technologies Co., Ltd. 1998-2014. All rights reserved.  
* File name: WSEC_Config.h
* Decription: 本文件需要APP程序员对编译选项进行简单的配置
*********************************************************************************/
#ifndef __WIRELESS_SEC_CONFIG_D4513A042DC_RF3F_4E427
#define __WIRELESS_SEC_CONFIG_D4513A042DC_RF3F_4E427

#ifdef __cplusplus
extern "C" {
#endif

/*================================================
       1. 编译开关
================================================*/
//#define WSEC_WIN32
//#define WSEC_LINUX
//#define WSEC_DEBUG
//#define WSEC_TRACE_MEMORY        /* CBB开发调测用, APP不理会 */
//#define WSEC_COMPILE_SDP         /* 启用SDP(敏感数据保护)子CBB */
//#define WSEC_COMPILE_CAC_IPSI    /* 启用基于iPSI的CAC(加密算法库适配)子CBB */
//#define WSEC_COMPILE_CAC_OPENSSL /* 启用基于OpenSSL的CAC(加密算法库适配)子CBB */

/*================================================
       2. CPU寻址模式
================================================*/
#define WSEC_CPU_ENDIAL_AUTO_CHK (0) /* 程序自动检测 */
#define WSEC_CPU_ENDIAL_BIG      (1) /* 大端对齐 */
#define WSEC_CPU_ENDIAL_LITTLE   (2) /* 小端对齐 */
#define WSEC_CPU_ENDIAN_MODE     WSEC_CPU_ENDIAL_AUTO_CHK /* 如果不定义该宏，则由程序自动检测 */

/*================================================
        3. 静态参数
================================================*/
#define WSEC_DOMAIN_NUM_MAX          (1024) /* Domain最大个数 */
#define WSEC_DOMAIN_KEY_TYPE_NUM_MAX (16)   /* 每个Domain拥有KeyType的最大个数 */
#define WSEC_MK_NUM_MAX              (4096) /* Keystore文件中允许存储MK最大数量 */
#define WSEC_ENABLE_BLOCK_MILSEC     (10)   /* 耗时操作, 允许连续占用CPU的时间(单位: 毫秒) */
//#define WSEC_ERR_CODE_BASE (0) /* 预留给CBB的起始错误码, APP必须显式定义 */

/*================================================
        4. 其它
================================================*/
//#define WSEC_WRI_LOG_AUTO_END_WITH_CRLF /* 如果产品记录日志时自动在尾部添加了回车换行则定义该宏，否则注释之 */

#ifdef __cplusplus
}
#endif  /* __cplusplus */

#endif/* __WIRELESS_SEC_CONFIG_D4513A042DC_RF3F_4E427 */
