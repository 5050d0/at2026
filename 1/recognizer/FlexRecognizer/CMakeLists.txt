add_library(FlexRecognizer FlexRecognizer.cpp FlexRecognizer.h lex.yy.cc lex.yy.h)
target_link_libraries(FlexRecognizer IRecognizer)
