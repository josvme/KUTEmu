
int main() {
    char *hello = "Hello C ";
    char *m = (char*) 0x10000000;
    for(int i=0; i < 8; i++) {
        *m = hello[i];
    }

    return 0;
}