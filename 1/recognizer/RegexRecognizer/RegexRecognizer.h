//
// Created by kirill on 2/17/26.
//
#pragma once

#include "recognizer/IRecognizer.h"
#include <unordered_map>
#include <regex>

class RegexRecognizer final : public IRecognizer {
    std::regex static const main_regex;
    static std::vector<std::string> const allowed_types;

    std::unordered_map<std::string, std::string> KnownVariables;

public:
    std::pair<bool, std::string> Recognize(std::string row) override;

    void reset() override;;
};
