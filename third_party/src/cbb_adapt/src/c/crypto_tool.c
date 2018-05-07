#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "aes_crypto.h"


int main(int argc,char *argv[])
{
    int iRet = 0;
    char *pPlainFile   = NULL;
    char *pEncryptFile = NULL;
    char *pTextString  = NULL;
    char *pEncryptString  = NULL;

    iRet = aesInit();
    if (iRet != 0)
    {
        printf("init error iRet=%d\n", iRet);
        return 1;
    }
    
    if (argc < 2 || 0 == strcmp(argv[1], "-help"))
    {
        printf("Usage: crypto_tool -file plainFileName encryptFileName\n"
               "or     crypto_tool -text textString\n");
    }
    else if (argc == 3 && 0 == strcmp(argv[1], "-text"))
    {
        pTextString   = argv[2];
        pEncryptString = aesEncrypt(0, pTextString);
        if (pEncryptString == NULL)
        {
            printf("File encrypt error, iRet=%d\n", iRet);
            return 1;
        }
        printf("%s\n", pEncryptString);
        free(pEncryptString);
        pEncryptString = NULL;
    }
    else if (argc == 4 && 0 == strcmp(argv[1], "-file"))
    {
        pPlainFile   = argv[2];
        pEncryptFile = argv[3];

        iRet = aesFileEncrypt(0, pPlainFile, pEncryptFile);
        if (iRet != 0)
        {
            printf("File encrypt error, iRet=%d\n", iRet);
            return 1;
        }
    }
    else
    {
        printf("Usage: crypto_tool -file plainFileName encryptFileName\n"
               "or     crypto_tool -text textString\n");
    }

    return 0;
}