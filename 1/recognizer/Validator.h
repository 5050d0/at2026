//
// Created by kirill on 3/3/26.
//

#pragma once
#include "RegexRecognizer/RegexRecognizer.h"

class Validator {
    std::unique_ptr<IRecognizer> recognizer;
    std::unordered_map<std::string, std::string> KnownVariables;

public:
    Validator() = default;

    explicit Validator(std::unique_ptr<RegexRecognizer> ptr);

    std::pair<bool, std::string> Validate(std::string line);

    void reset();
};
