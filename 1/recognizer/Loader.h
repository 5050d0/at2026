//
// Created by kirill on 2/16/26.
//

#pragma once

#include <map>
#include <memory>
#include <string>
#include <vector>

#include "IRecognizer.h"
#include "RegexRecognizer/RegexRecognizer.h"

class Loader {
    std::map<std::string, std::unique_ptr<IRecognizer> > recognizers;

public:
    Loader() {
        recognizers.emplace("regex", std::make_unique<RegexRecognizer>());
    }

    std::map<std::string, std::unique_ptr<IRecognizer> > &get_recognizers() {
        return recognizers;
    }
};


