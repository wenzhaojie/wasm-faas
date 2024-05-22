#include <wasmedge/wasmedge.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>

#define DATA_SIZE_IN_MB 1 // 定义特定大小的数据大小，单位为MB

#define DATA_SIZE (DATA_SIZE_IN_MB * 1024 * 1024) // 将数据大小转换为字节

void generate_data(char *data, size_t size) {
    // 生成特定大小的字符串
    for (size_t i = 0; i < size - 1; ++i) {
        data[i] = 'A' + (rand() % 26);
    }
    data[size - 1] = '\0'; // 添加字符串结束符
}

double current_timestamp() {
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (double)tv.tv_sec + (double)tv.tv_usec / 1000000;
}

int main(int argc, const char *argv[]) {
    if (argc != 2) {
        printf("Usage: %s <wasm_file>\n", argv[0]);
        return 1;
    }

    /* Create the configure context */
    WasmEdge_ConfigureContext *conf_cxt = WasmEdge_ConfigureCreate();
    /* The configure and store context to the VM creation can be NULL. */
    WasmEdge_VMContext *vm_cxt = WasmEdge_VMCreate(conf_cxt, NULL);

    /* The parameters and returns arrays. */
    WasmEdge_Value params[1];
    /* Function name. */
    WasmEdge_String func_name = WasmEdge_StringCreateByCString("bindwidth");

    /* 准备数据 */
    char *data = (char *)malloc(DATA_SIZE);
    generate_data(data, DATA_SIZE);
    params[0] = WasmEdge_ValueGenI32(DATA_SIZE);
    /* 调用Wasm函数并计时 */
    double start_time = current_timestamp();
    /* Run the WASM function from file. */
    WasmEdge_Value returns[1];
    WasmEdge_Result res = WasmEdge_VMRunWasmFromFile(vm_cxt, argv[1], func_name, params, 1, returns, 1);
    double end_time = current_timestamp();

    if (WasmEdge_ResultOK(res)) {
        double elapsed_time = end_time - start_time;
        double bandwidth = (DATA_SIZE / (1024 * 1024)) / elapsed_time; // 将带宽转换为MB/s

        /* 打印结果 */
        printf("Data Size: %d MB\n", DATA_SIZE_IN_MB);
        printf("Elapsed Time: %.6f seconds\n", elapsed_time);
        printf("Bandwidth: %.2f MB/second\n", bandwidth);
    } else {
        printf("Error message: %s\n", WasmEdge_ResultGetMessage(res));
    }

    /* 释放资源 */
    WasmEdge_VMDelete(vm_cxt);
    WasmEdge_ConfigureDelete(conf_cxt);
    WasmEdge_StringDelete(func_name);
    free(data);

    return 0;
}
