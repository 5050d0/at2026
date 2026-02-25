//
// Created by kirill on 2/17/26.
//

#pragma once
#include <string>

class IRecognizer {
public:
    virtual ~IRecognizer() = default;
    virtual std::pair<bool, std::string> Recognize(std::string row) = 0;
    virtual void reset() = 0;
};
