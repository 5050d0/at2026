//
// Created by kirill on 2/17/26.
//
#pragma once

#include "recognizer/IRecognizer.h"

#include <boost/regex.hpp>

class RegexRecognizer final : public IRecognizer {
    boost::regex static const main_regex;
    static std::vector<std::string> const allowed_types;

public:
    std::optional<RecResult> Recognize(std::string row) override;

    void reset() override;;
};
