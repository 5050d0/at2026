//
// Created by kirill on 2/17/26.
//
#pragma once

#include <unordered_map>
#include "recognizer/IRecognizer.h"

class RegexRecognizer final : public IRecognizer {
    std::unordered_map<std::string, std::string> KnownVariables;

public:
    std::pair<bool, std::string> Recognize(std::string row) override;

    void reset() override;;
};
