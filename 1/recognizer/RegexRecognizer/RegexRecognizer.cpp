//
// Created by kirill on 2/17/26.
//

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

boost::regex const RegexRecognizer::main_regex{
    "^ *(int|short|long) +([a-zA-Z][a-zA-Z0-9]{0,15}) *= *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+))(?: *([%/*]) *(?:([a-zA-Z][a-zA-Z0-9]{0,15})|([0-9]+)))? *; *$"
};
std::vector<std::string> const RegexRecognizer::allowed_types = {
    "int", "short", "long"
};

std::optional<RecResult> RegexRecognizer::Recognize(const std::string &row) {
    boost::smatch match;

    if (!boost::regex_search(row, match, main_regex)) {
        return std::nullopt;
    }
    return RecResult{
        .vartype = match[1].str(), .lvar = match[2].str(), .rvar1 = match[3].str(), .rvar2 = match[6].str()
    };
}

void RegexRecognizer::reset() {
    *this = {};
}
