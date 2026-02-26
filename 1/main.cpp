#include <iostream>
#include <regex>

#include "recognizer/Loader.h"

using namespace std;

int main() {
    Loader loader;
    loader.get_recognizers();
    loader.get_recognizers().begin()->second->Recognize("");
}

