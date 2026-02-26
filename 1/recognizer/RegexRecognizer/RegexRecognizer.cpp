//
// Created by kirill on 2/17/26.
//

#include <regex>

#include "RegexRecognizer.h"

#include <iostream>


// enum class CaptureGroups: char {
//     CAPTURE_FULL = 0,
//     CAPTURE_LEFT_TYPE = 1,
//     CAPTURE_LEFT_NAME = 2,
//     CAPTURE_OP1_NAME = 3,
//     CAPTURE_OP1_NUM = 4,
//     CAPTURE_OP = 5,
//     CAPTURE_OP2_NAME = 6,
//     CAPTURE_OP2_NUM = 7
// };

std::regex const RegexRecognizer::main_regex{
    R"(^(int|short|long) +([a-zA-Z][a-zA-Z0-9]{0,15}) *= *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+))(?: *([%/*]) *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+)))?;$)"
};
std::vector<std::string> const RegexRecognizer::allowed_types = {
    "int", "short", "long"
};

std::pair<bool, std::string> RegexRecognizer::Recognize(std::string row) {
    std::smatch match;
    if (!std::regex_search(row, match, main_regex)) {
        return {false, ""};
    }
    if (std::ranges::find(allowed_types, match[1]) == std::ranges::end(allowed_types)) {
        return {false, ""};
    }

    auto found = KnownVariables.find(match[2]);
    if (found == KnownVariables.end()) {
        KnownVariables[match[2]] = match[1];
    } else {
        if (found->second != match[1]) {
            return {
                false,
                std::format("Redeclaration of variable {} with type {} (was type {})", match[2].str(), match[1].str(),
                            found->second)
            };
        }
    }
    // todo надо ли проверять чтобы используемые переменные совпадали типом?

    // size_t j = 0;
    // for (auto i: match) {
    //     std::cout << j++ << " : " << i.str() << std::endl;
    // }


    return {true, ""};
}

void RegexRecognizer::reset() {
    *this = {};
}
