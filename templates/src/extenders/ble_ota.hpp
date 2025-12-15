#pragma once

#include <zephyr/mgmt/mcumgr/mgmt/callbacks.h>

namespace zbhome
{
    namespace experimental
    {
        class BLEOTA
        {
        public:
            static int on_main();

        private:
            static mgmt_cb_return on_img_uploaded(uint32_t event, enum mgmt_cb_return prev_status,
                                                  int32_t *rc, uint16_t *group, bool *abort_more, void *data,
                                                  size_t data_size);
            static mgmt_cb_return on_reset_evt(uint32_t event, enum mgmt_cb_return prev_status,
                                               int32_t *rc, uint16_t *group, bool *abort_more, void *data,
                                               size_t data_size);
        };
    }
}
