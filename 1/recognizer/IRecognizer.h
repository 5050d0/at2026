//
// Created by kirill on 2/17/26.
//

#pragma once
#include <optional>
#include <string>
#include <vector>

struct RecResult {
    std::string vartype;
    std::string lvar;
    std::string rvar1;
    std::string rvar2;
};

class IRecognizer {
public:
    virtual ~IRecognizer() = default;

    virtual std::optional<RecResult> Recognize(std::string row) = 0;

    virtual void reset() = 0;
};
