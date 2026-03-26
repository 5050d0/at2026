#include <iostream>
#include <regex>

#include "recognizer/Loader.h"
#include "recognizer/FlexRecognizer/FlexRecognizer.h"

using namespace std;

int main() {
    Loader loader;
    auto &a = loader.get_recognizers()["flex"];
    auto res = a.Validate("int Dsasf323 = ab /     5;");


    res = a.Validate("long ab = 4;");
    // FlexRecognizer flexRecognizer;
    // auto res = flexRecognizer.Recognize("int ab= 444;");
    // res = flexRecognizer.Recognize("long ab = 4;");
    // reurn 0;
}

