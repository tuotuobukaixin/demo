/*******************************************************************************
* Copyright @ Huawei Technologies Co., Ltd. 1998-2014. All rights reserved.  
* File name: WSEC_ErrorCode.h
* Decription: 错误码
*********************************************************************************/
#ifndef __WIRELESS_ERROR_CODE_D413A42DCRF_3F4E_427
#define __WIRELESS_ERROR_CODE_D413A42DCRF_3F4E_427

#include "WSEC_Config.h"

#ifdef __cplusplus
extern "C" {
#endif

/* 致CBB开发者: 在IDE Project定义该宏, 不要在'WSEC_Config.h'定义, 以促使APP给出正确的定义 */
#ifndef WSEC_ERR_CODE_BASE 
    #error Please define the 'WSEC_ERR_CODE_BASE' in 'WSEC_Config.h'.
#endif

#define WSEC_ERROR_CODE(seq) ((WSEC_ERR_T)(WSEC_ERR_CODE_BASE + seq))

#define WSEC_SUCCESS                                               (WSEC_ERR_T)0    /* 成功 */
#define WSEC_FAILURE                                           WSEC_ERROR_CODE(1)   /* 通用错误 */

/* 文件操作错误 */
#define WSEC_ERR_OPEN_FILE_FAIL                                WSEC_ERROR_CODE(11)   /* 打开文件失败 */
#define WSEC_ERR_READ_FILE_FAIL                                WSEC_ERROR_CODE(12)   /* 读文件失败 */
#define WSEC_ERR_WRI_FILE_FAIL                                 WSEC_ERROR_CODE(13)   /* 写文件失败 */
#define WSEC_ERR_GET_FILE_LEN_FAIL                             WSEC_ERROR_CODE(14)   /* 获取文件长度失败 */
#define WSEC_ERR_FILE_FORMAT                                   WSEC_ERROR_CODE(15)   /* 文件格式错误 */

/* 内存操作错误 */
#define WSEC_ERR_MALLOC_FAIL                                   WSEC_ERROR_CODE(51)  /* 内存分配失败 */
#define WSEC_ERR_MEMCPY_FAIL                                   WSEC_ERROR_CODE(52)  /* 内存拷贝失败 */
#define WSEC_ERR_MEMCLONE_FAIL                                 WSEC_ERROR_CODE(53)  /* 内存克隆失败 */
#define WSEC_ERR_STRCPY_FAIL                                   WSEC_ERROR_CODE(54)  /* 字符串拷贝失败 */
#define WSEC_ERR_OPER_ARRAY_FAIL                               WSEC_ERROR_CODE(55)  /* 数组操作失败 */

/* 安全函数处理错误 */
#define WSEC_ERR_CRPTO_LIB_FAIL                                WSEC_ERROR_CODE(101) /* 安全函数库(iPSI)操作失败 */
#define WSEC_ERR_GEN_HASH_CODE_FAIL                            WSEC_ERROR_CODE(102) /* 生成Hash值失败 */
#define WSEC_ERR_HASH_NOT_MATCH                                WSEC_ERROR_CODE(103) /* Hash值不匹配 */
#define WSEC_ERR_INTEGRITY_FAIL                                WSEC_ERROR_CODE(104) /* 完整性被破坏 */
#define WSEC_ERR_HMAC_FAIL                                     WSEC_ERROR_CODE(105) /* HMAC失败 */
#define WSEC_ERR_HMAC_AUTH_FAIL                                WSEC_ERROR_CODE(106) /* HMAC验证失败 */
#define WSEC_ERR_GET_RAND_FAIL                                 WSEC_ERROR_CODE(107) /* 获取随机数失败 */
#define WSEC_ERR_PBKDF2_FAIL                                   WSEC_ERROR_CODE(108) /* 推演密钥失败 */
#define WSEC_ERR_ENCRPT_FAIL                                   WSEC_ERROR_CODE(109) /* 数据加密失败 */
#define WSEC_ERR_DECRPT_FAIL                                   WSEC_ERROR_CODE(110) /* 数据解密失败 */
#define WSEC_ERR_GET_ALG_NAME_FAIL                             WSEC_ERROR_CODE(111) /* 获取安全算法名失败 */

/* 函数调用类错误 */
#define WSEC_ERR_INVALID_ARG                                   WSEC_ERROR_CODE(151) /* 非法参数 */
#define WSEC_ERR_OUTPUT_BUFF_NOT_ENOUGH                        WSEC_ERROR_CODE(152) /* 输出缓冲区不足 */
#define WSEC_ERR_INPUT_BUFF_NOT_ENOUGH                         WSEC_ERROR_CODE(153) /* 输入缓冲区不足 */
#define WSEC_ERR_CANCEL_BY_APP                                 WSEC_ERROR_CODE(154) /* APP取消操作 */
#define WSEC_ERR_INVALID_CALL_SEQ                              WSEC_ERROR_CODE(155) /* APP调用顺序错误 */

/* 系统操作错误 */
#define WSEC_ERR_GET_CURRENT_TIME_FAIL                         WSEC_ERROR_CODE(201) /* 获取当前时间失败 */

/* KMC错误 */
#define WSEC_ERR_KMC_CALLBACK_KMCCFG_FAIL                      WSEC_ERROR_CODE(251) /* 回调获取KMC配置数据失败 */
#define WSEC_ERR_KMC_KMCCFG_INVALID                            WSEC_ERROR_CODE(252) /* KMC配置数据非法 */
#define WSEC_ERR_KMC_KSF_DATA_INVALID                          WSEC_ERROR_CODE(253) /* Keystore存在非法配置数据 */
#define WSEC_ERR_KMC_INI_MUL_CALL                              WSEC_ERROR_CODE(254) /* 多次调用初始化 */
#define WSEC_ERR_KMC_NOT_KSF_FORMAT                            WSEC_ERROR_CODE(255) /* 不是Keystore文件格式 */
#define WSEC_ERR_KMC_READ_DIFF_VER_KSF_FAIL                    WSEC_ERROR_CODE(256) /* 读取其它版本的Keystore文件失败 */
#define WSEC_ERR_KMC_READ_MK_FAIL                              WSEC_ERROR_CODE(257) /* 读取MK失败 */
#define WSEC_ERR_KMC_MK_LEN_TOO_LONG                           WSEC_ERROR_CODE(258) /* MK密钥超长 */
#define WSEC_ERR_KMC_REG_REPEAT_MK                             WSEC_ERROR_CODE(259) /* 试图注册重复的MK */
#define WSEC_ERR_KMC_ADD_REPEAT_DOMAIN                         WSEC_ERROR_CODE(260) /* 试图增加重复的Domain(ID重复) */
#define WSEC_ERR_KMC_ADD_REPEAT_KEY_TYPE                       WSEC_ERROR_CODE(261) /* 试图增加重复的KeyType(同一Domain下KeyType重复) */
#define WSEC_ERR_KMC_ADD_REPEAT_MK                             WSEC_ERROR_CODE(262) /* 试图增加重复的MK(同一Domain下KeyId重复) */
#define WSEC_ERR_KMC_DOMAIN_MISS                               WSEC_ERROR_CODE(263) /* DOMAIN不存在 */
#define WSEC_ERR_KMC_DOMAIN_KEYTYPE_MISS                       WSEC_ERROR_CODE(264) /* DOMAIN KeyType不存在 */
#define WSEC_ERR_KMC_DOMAIN_NUM_OVERFLOW                       WSEC_ERROR_CODE(265) /* DOMAIN配置数量超限 */
#define WSEC_ERR_KMC_KEYTYPE_NUM_OVERFLOW                      WSEC_ERROR_CODE(266) /* KeyType配置数量超限 */
#define WSEC_ERR_KMC_MK_NUM_OVERFLOW                           WSEC_ERROR_CODE(267) /* MK数量超限 */
#define WSEC_ERR_KMC_MK_MISS                                   WSEC_ERROR_CODE(268) /* MK不存在 */
#define WSEC_ERR_KMC_RECREATE_MK                               WSEC_ERROR_CODE(269) /* 重新创建MK失败 */
#define WSEC_ERR_KMC_CBB_NOT_INIT                              WSEC_ERROR_CODE(270) /* CBB尚未初始化 */
#define WSEC_ERR_KMC_CANNOT_REG_AUTO_KEY                       WSEC_ERROR_CODE(271) /* 不能注册系统自动生成的密钥 */
#define WSEC_ERR_KMC_CANNOT_RMV_ACTIVE_MK                      WSEC_ERROR_CODE(272) /* 不能删除处于活动状态的MK */
#define WSEC_ERR_KMC_CANNOT_SET_EXPIRETIME_FOR_INACTIVE_MK     WSEC_ERROR_CODE(273) /* 不能对inactive的MK设置过期时间 */
#define WSEC_ERR_KMC_RK_GENTYPE_REJECT_THE_OPER                WSEC_ERROR_CODE(274) /* RK的生成方式不支持该操作 */
#define WSEC_ERR_KMC_MK_GENTYPE_REJECT_THE_OPER                WSEC_ERROR_CODE(275) /* MK的生成方式不支持该操作 */
#define WSEC_ERR_KMC_ADD_DOMAIN_DISCREPANCY_MK                 WSEC_ERROR_CODE(276) /* 待增DOMAIN与残留的MK矛盾 */
#define WSEC_ERR_KMC_IMPORT_MK_CONFLICT_DOMAIN                 WSEC_ERROR_CODE(277) /* 导入的MK与Domain配置冲突 */
#define WSEC_ERR_KMC_CANNOT_ACCESS_PRI_DOMAIN                  WSEC_ERROR_CODE(278) /* 不能访问CBB私有Domain */

/* SDP错误 */
#define WSEC_ERR_SDP_PWD_VERIFY_FAIL                           WSEC_ERROR_CODE(351) /* 密码密文校验失败 */
#define WSEC_ERR_SDP_CONFIG_INCONSISTENT_WITH_USE              WSEC_ERROR_CODE(352) /* 配置数据与使用不一致 */
#define WSEC_ERR_SDP_INVALID_CIPHER_TEXT                       WSEC_ERROR_CODE(353) /* 密文格式解析错误 */
#define WSEC_ERR_SDP_VERSION_INCOMPATIBLE                      WSEC_ERROR_CODE(354) /* 密文版本与当前版本不兼容 */
#define WSEC_ERR_SDP_ALG_NOT_SUPPORTED                         WSEC_ERROR_CODE(355) /* 算法不存在或不支持 */
#define WSEC_ERR_SDP_DOMAIN_UNEXPECTED			               WSEC_ERROR_CODE(356) /* 密文来自非预期Domain */

#define WSEC_ERR_MAX                                           WSEC_ERROR_CODE(5000) /* CBB最大错误码 */

#ifdef __cplusplus
}
#endif  /* __cplusplus */

#endif/* __WIRELESS_ERROR_CODE_D413A42DCRF_3F4E_427 */
