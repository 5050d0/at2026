#include <iostream>
#include <regex>

#include "recognizer/Loader.h"

using namespace std;

int main() {
    Loader loader;
    loader.get_recognizers();
    loader.get_recognizers().begin()->second->Recognize("");


    std::vector<std::string> const tst = {
        "int ab= frgfr;",
        "short Dsasf323 = ab +     5;",
        "long vA2r23=5+6;",
        "a ad=5;", "",
        "int 3d = ab+Dsasf323;",
        "int a =;"
    };
}

