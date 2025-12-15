// extenders.hpp
#pragma once

#include <functional>
#include <vector>

namespace zbhome
{
    namespace experimental
    {
        namespace extenders
        {
            class Registry
            {
            public:
                using Action = std::function<int()>;

                static Registry &instance()
                {
                    static Registry theInstance;
                    return theInstance;
                }

                void register_extender(const Action &a)
                {
                    m_extenders.push_back(a);
                }

                void run_all()
                {
                    for (auto &extender : m_extenders)
                        extender();
                }

            private:
                std::vector<Action> m_extenders;
            };
        }
    }
}

#define REGISTER_EXTENDER(name, action)                                 \
    namespace zbhome                                                    \
    {                                                                   \
        namespace experimental                                          \
        {                                                               \
            namespace extenders                                         \
            {                                                           \
                struct name##_registrar                                 \
                {                                                       \
                    name##_registrar()                                  \
                    {                                                   \
                        Registry::instance().register_extender(action); \
                    }                                                   \
                };                                                      \
                static name##_registrar name##_registrar_instance;      \
            }                                                           \
        }                                                               \
    }