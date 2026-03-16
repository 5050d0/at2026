//
// Created by kirill on 2/16/26.
//

#pragma once

#include <map>
#include <memory>
#include <string>
#include <vector>

#include "IRecognizer.h"
#include "Validator.h"

#include "FlexRecognizer/FlexRecognizer.h"
#include "RegexRecognizer/RegexRecognizer.h"
#include "SMCRecognizer/SmcRecognizer.h"

class Loader {
    std::map<std::string, Validator> recognizers;

public:
    Loader() {
        recognizers.emplace("regex", Validator{std::make_unique<RegexRecognizer>()});
        recognizers.emplace("smc", Validator{std::make_unique<SmcRecognizer>()});
        recognizers.emplace("flex", Validator{std::make_unique<FlexRecognizer>()});
    }

    std::map<std::string, Validator> &get_recognizers() {
        return recognizers;
    }
};


