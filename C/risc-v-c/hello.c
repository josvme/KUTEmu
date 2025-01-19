
int main() {
    char *hello = "Hello";
    char *m = (char*) 0x10000000;
    for(int i=0; i < 6; i++) {
        *m = hello[i];
    }
    return 0;
}