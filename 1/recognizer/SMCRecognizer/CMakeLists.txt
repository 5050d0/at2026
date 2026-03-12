add_library(SmcRecognizer SmcRecognizer.cpp SmcRecognizer.h)


add_library(Smachine Smachine.h Smachine.cpp recognizer_sm.cpp recognizer_sm.h)

include(FetchContent)

FetchContent_Declare(
        smc
        GIT_REPOSITORY https://git.code.sf.net/p/smc/git-code
        GIT_TAG master
)

FetchContent_MakeAvailable(smc)

add_library(smc_statemap INTERFACE)

target_include_directories(smc_statemap
        INTERFACE
        ${smc_SOURCE_DIR}/lib/C++
)


add_library(smc::statemap ALIAS smc_statemap)


target_link_libraries(Smachine
        PUBLIC smc_statemap
)

target_link_libraries(SmcRecognizer
        PUBLIC Smachine
        PUBLIC IRecognizer
)
