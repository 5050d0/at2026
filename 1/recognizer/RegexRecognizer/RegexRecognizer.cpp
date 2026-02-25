//
// Created by kirill on 2/17/26.
//

#include <regex>

#include "RegexRecognizer.h"

#include <iostream>


enum class CaptureGroups: char {
    CAPTURE_FULL = 0,
    CAPTURE_LEFT_TYPE = 1,
    CAPTURE_LEFT_NAME = 2,
    CAPTURE_OP1_NAME = 3,
    CAPTURE_OP1_NUM = 4,
    CAPTURE_OP = 5,
    CAPTURE_OP2_NAME = 6,
    CAPTURE_OP2_NUM = 7
};

std::regex const RegexRecognizer::main_regex{
    R"(^(int|short|long) +([a-zA-Z][a-zA-Z0-9]{0,15}) *= *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+))(?: *([%*/+-]) *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+)))?;$)"
};

std::pair<bool, std::string> RegexRecognizer::Recognize(std::string row) {
    for (const auto &str: tst) {
        std::smatch match;
        if (!std::regex_search(str, match, main_regex)) {
            return {false, ""};
        }
        size_t j = 0;
        for (auto i: match) {
            std::cout << j++ << " : " << i.str() << std::endl;
        }
        // for (auto tok_begin = std::sregex_iterator(str.begin(), str.end(), main_regex);
        //      tok_begin != std::sregex_iterator{}; ++tok_begin) {
        //     std::cout << tok_begin->str() << std::endl;
        // }
    }

    // std::regex lvar_name_regex
    //^(int|short|long) +[a-zA-Z]+[a-zA-Z0-9]* *=
    //KnownVariables[]
    return {true, ""};
}

void RegexRecognizer::reset() {
    *this = {};
}
