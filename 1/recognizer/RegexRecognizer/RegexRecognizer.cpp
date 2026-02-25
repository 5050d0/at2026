//
// Created by kirill on 2/17/26.
//

#include <regex>

#include "RegexRecognizer.h"

#include <iostream>


std::pair<bool, std::string> RegexRecognizer::Recognize(std::string row) {
    std::regex main_regex{
        "^(int|short|long) +([a-zA-Z][a-zA-Z0-9]{0,15}) *= *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+))(?: *([%*/+-]) *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+)))?;$"
    };
    std::vector<std::string> tst = {
        "int ab= frgfr;",
        "short Dsasf323 = ab +     5;",
        "long vA2r23=5+6;",
        "a ad=5;", "",
        "int 3d = ab+Dsasf323;",
        "int a =;"
    };
    for (auto &str: tst) {
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
