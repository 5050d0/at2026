//
// Created by kirill on 3/3/26.
//

#include "Validator.h"

Validator::Validator(std::unique_ptr<RegexRecognizer> ptr) : recognizer(std::move(ptr)) {
}

std::pair<bool, std::string> Validator::Validate(std::string line) {
    const auto res = recognizer->Recognize(line);

    if (!res.has_value()) return {false, ""};


    if (!((res->rvar1.empty() || KnownVariables.contains(res->rvar1)) && (
              res->rvar2.empty() || KnownVariables.contains(res->rvar2)))) {
        return {false, ""};
    }

    const auto &f = KnownVariables.find(res.value().lvar);
    if (f == KnownVariables.end()) {
        KnownVariables[res.value().lvar] = res.value().vartype;
        return {true, ""};
    }
    if (f->second == res.value().vartype) {
        return {true, ""};
    }
    return {
        true,
        std::format("Variable {} redeclared with type {} (was {})", res.value().lvar, res.value().vartype, f->second)
    };
}

void Validator::reset() {
    recognizer->reset();
    KnownVariables.clear();
}
