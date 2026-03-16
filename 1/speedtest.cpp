

#include <iostream>
#include "recognizer/Loader.h"

using namespace std;

int main() {
    Loader loader;
    auto &a = loader.get_recognizers()["flex"];
    auto res = a.Validate("short ab = 1;");

    res = a.Validate("int ab= 444;");


    res = a.Validate("long ab = 4;");
}

