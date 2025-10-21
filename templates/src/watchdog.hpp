#pragma once

/*
 * Copyright (c) 2015 Intel Corporation
 * Copyright (c) 2018 Nordic Semiconductor
 * Copyright (c) 2019 Centaur Analytics, Inc
 * Copyright 2023 NXP
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// This header will allow enabling watchdog,
// which in turn would reset SoC in case
// of some lock up.

#include <zephyr/kernel.h>
#include <zephyr/device.h>
#include <zephyr/drivers/watchdog.h>

#define WDT_MAX_WINDOW 65000U
#define WDG_FEED_INTERVAL 60000U

// This check is present as we don't
// force watchdog on in the devicetree.
#if DT_HAS_COMPAT_STATUS_OKAY(nordic_nrf_wdt)
const struct device *const wdt = DEVICE_DT_GET(DT_ALIAS(watchdog0));
#else
const struct device *const wdt = NULL;
#endif

static struct wdt_window wdt_window_cfg = {
    .min = WDT_MIN_WINDOW,
    .max = WDT_MAX_WINDOW,
};

static struct wdt_timeout_cfg wdt_config = {
    /* Expire watchdog after max window */
    .window = wdt_window_cfg,
    /* Reset SoC when watchdog timer expires. */
    .flags = WDT_FLAG_RESET_SOC,
};

static int wdt_channel_id;

void feed_watchdog(struct k_timer *dummy)
{
    wdt_feed(wdt, wdt_channel_id);
};

K_TIMER_DEFINE(wdt_timer, feed_watchdog, NULL);

int setup_watchdog()
{
    // If we don't have watchdog - don't sweat it, for now.
    if (!wdt)
    {
        return -1;
    }

    int err;
    if (!device_is_ready(wdt))
    {
        return -2;
    }

    wdt_channel_id = wdt_install_timeout(wdt, &wdt_config);
    if (wdt_channel_id < 0)
    {
        return -3;
    }

    err = wdt_setup(wdt, WDT_OPT_PAUSE_HALTED_BY_DBG);
    if (err < 0)
    {
        return -4;
    }

    // Schedule feeding of watchdog
    k_timer_start(&wdt_timer, K_SECONDS(WDG_FEED_INTERVAL / 1000), K_SECONDS(WDG_FEED_INTERVAL / 1000));

    return 0;
};
