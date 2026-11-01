import {getColorSync, getPaletteSync} from "colorthief";
import chroma from "chroma-js";

// ======================== 常量定义 ======================== //
const RGB_MAX_VALUE = 255;
const DEFAULT_COLOR_COUNT = 8;
const WHITE_TEXT_LUMINANCE_THRESHOLD = 0.55; // 判断是否使用白色文字的亮度阈值

// 对比度常量
const MIN_CONTRAST = 4.5; // 文本与背景的最小对比度
const MAX_CONTRAST_ATTEMPTS = 5; // 调整对比度的最大尝试次数
const MIN_BASE_CONTRAST = 2.5; // 主色调与基础色的最小对比度
const MIN_SUBTITLE_CONTRAST_TEXT = 3.2; // 副标题与文本的最小对比度
const MIN_SUBTITLE_CONTRAST_BG = 4.0; // 副标题与背景的最小对比度
const MAX_SUBTITLE_CONTRAST_ATTEMPTS = 6; // 调整副标题对比度的最大尝试次数

// 图片加载常量
const IMAGE_LOAD_TIMEOUT = 30000; // 30秒超时
const DEFAULT_BACKGROUND_IMAGE = "url('/images/default.jpg')";

// 色相偏移常量
const HUE_SHIFT_200_250 = 10;
const HUE_SHIFT_250_320 = 20;
const HUE_SHIFT_DEFAULT = 25;
const HUE_MAX = 360;

// 饱和度/亮度调整常量 (selectedColor)
const SATURATION_ADD = 0.22;
const SATURATION_MAX = 0.88;
const LUMINANCE_MIN = 0.48;
const LUMINANCE_MAX = 0.68;
const BRIGHTEN_ADD = 0.18;
const SATURATE_ADD = 0.28;

// 黄色调优化常量
const YELLOW_HUE_MIN = 40;
const YELLOW_HUE_MAX = 65;
const YELLOW_LUMINANCE_THRESHOLD = 0.7;
const YELLOW_DARKEN_AMOUNT = 0.5;
const YELLOW_HUE_SHIFT_AMOUNT = 20;

// 混合颜色常量 (backgroundBlendColor)
const MIX_BLACK_WHITE_RATIO_BLEND = 0.28;
const LUMINANCE_MULTIPLIER_BLEND = 1.08;
const DESATURATE_AMOUNT_BLEND = 0.25;
const ALPHA_WHITE_TEXT_BLEND = 0.25;
const ALPHA_BLACK_TEXT_BLEND = 0.12;

// 混合颜色常量 (backgroundRightColor)
const MIX_BASE_SELECTED_RATIO_RIGHT = 0.3;
const SATURATION_MULTIPLIER_RIGHT = 1.2;
const LUMINANCE_MULTIPLIER_RIGHT = 1.1;
const ALPHA_WHITE_TEXT_RIGHT = 0.2;
const ALPHA_BLACK_TEXT_RIGHT = 0.4;

// 副标题颜色常量
const SUBTITLE_BRIGHTEN_DARKEN_AMOUNT = 0.3;
const SUBTITLE_SATURATE_AMOUNT = 0.2;
const SUBTITLE_ALPHA = 0.95;
const SUBTITLE_HUE_SHIFT_WHITE = 10;
const SUBTITLE_HUE_SHIFT_BLACK = 15;

// body-blur-color 常量
const BODY_BLUR_BLACK_ALPHA = 0.16;
const BODY_BLUR_WHITE_ALPHA = 0.10;
const BODY_BLUR_MIX_RATIO = 0.3;
const BODY_BLUR_LUMINANCE_MULTIPLIER_WHITE = 0.85;
const BODY_BLUR_LUMINANCE_MULTIPLIER_BLACK = 1.15;


// ======================== 工具函数 ======================== //

/**
 * 计算颜色的感知亮度。
 * @param rgb 包含红、绿、蓝分量的数组。
 * @returns 感知亮度值（0-1）。
 */
const getPerceivedLuminance = (rgb: number[]): number => {
    // 根据 ITU-R BT.709 标准计算感知亮度
    return (
        0.299 * rgb[0] + 0.587 * rgb[1] + 0.114 * rgb[2]
    ) / RGB_MAX_VALUE;
};

/**
 * 提取图片主色调调色板并计算平均亮度。
 * @param img HTMLImageElement 对象。
 * @param colorCount 要提取的颜色数量。
 * @returns 平均亮度值。
 */
const getAverageLuminance = (img: HTMLImageElement, colorCount: number = DEFAULT_COLOR_COUNT): number => {
    const palette = getPaletteSync(img, {colorCount}).map((color) => color.array());
    const totalLuminance = palette.reduce((sum, color) => sum + getPerceivedLuminance(color), 0);
    return totalLuminance / palette.length;
};

/**
 * 判断是否应该使用白色文本以获得更好的可读性。
 * @param img HTMLImageElement 对象。
 * @returns 如果应该使用白色文本，则为 true；否则为 false。
 */
const shouldUseWhiteText = (img: HTMLImageElement): boolean =>
    getAverageLuminance(img) < WHITE_TEXT_LUMINANCE_THRESHOLD;

/**
 * 设置CSS变量。
 * @param vars 包含CSS变量名和值的对象。
 */
const setCSSVariables = (vars: Record<string, string>) => {
    const root = document.documentElement;
    Object.entries(vars).forEach(([key, value]) => {
        root.style.setProperty(`--${key}`, value);
    });
};

/**
 * 计算两种RGB颜色之间的对比度。
 * @param rgb1 第一个颜色的RGB数组。
 * @param rgb2 第二个颜色的RGB数组。
 * @returns 对比度值。
 */
const getContrast = (rgb1: number[], rgb2: number[]): number => {
    const luminance = (rgb: number[]) => {
        const a = rgb.map((v) => {
            v /= RGB_MAX_VALUE;
            return v <= 0.03928
                ? v / 12.92
                : Math.pow((v + 0.055) / 1.055, 2.4);
        });
        return 0.2126 * a[0] + 0.7152 * a[1] + 0.0722 * a[2];
    };
    const L1 = luminance(rgb1);
    const L2 = luminance(rgb2);
    return (Math.max(L1, L2) + 0.05) / (Math.min(L1, L2) + 0.05);
};

/**
 * 根据色相值计算智能偏移后的色相。
 * @param h 原始色相值。
 * @returns 偏移后的色相值。
 */
const calculateHueShift = (h: number): number => {
    if (h >= 200 && h <= 250) {
        return h - HUE_SHIFT_200_250;
    } else if (h > 250 && h < 320) {
        return h - HUE_SHIFT_250_320;
    } else {
        return (h + HUE_SHIFT_DEFAULT) % HUE_MAX;
    }
};

/**
 * 调整和优化主选颜色。
 * @param baseColor 基础颜色（从图片提取）。
 * @param h 基础色的色相。
 * @param s 基础色的饱和度。
 * @param l 基础色的亮度。
 * @returns 调整后的主选颜色。
 */
const adjustSelectedColor = (
    baseColor: chroma.Color,
    h: number,
    s: number,
    l: number
): chroma.Color => {
    // 动态调整饱和度和亮度，考虑背景亮度
    const dynamicSaturation = Math.min(s + SATURATION_ADD + (l < 0.5 ? 0.05 : -0.05), SATURATION_MAX); // 如果背景亮，饱和度稍微低一点，反之高一点
    const dynamicLuminance = Math.min(Math.max(l, LUMINANCE_MIN), LUMINANCE_MAX);

    let selectedColor = chroma.hsl(
        calculateHueShift(h),
        dynamicSaturation,
        dynamicLuminance
    )
        .saturate(SATURATE_ADD)
        .brighten(BRIGHTEN_ADD)
        .alpha(0.88);

    // 避免偏亮黄色
    if (h > YELLOW_HUE_MIN && h < YELLOW_HUE_MAX && l > YELLOW_LUMINANCE_THRESHOLD) {
        selectedColor = selectedColor.darken(YELLOW_DARKEN_AMOUNT).set("hsl.h", (h + YELLOW_HUE_SHIFT_AMOUNT) % HUE_MAX);
    }
    return selectedColor;
};

/**
 * 确保给定颜色与文本颜色之间有足够的对比度。
 * @param color 要调整的颜色。
 * @param textRGB 文本的RGB颜色数组。
 * @param useWhiteText 是否使用白色文本。
 * @returns 调整后的颜色。
 */
const ensureContrastWithText = (
    color: chroma.Color,
    textRGB: number[],
    useWhiteText: boolean
): chroma.Color => {
    let adjustedColor = color;
    let attempts = 0;
    while (getContrast(adjustedColor.rgb(), textRGB) < MIN_CONTRAST && attempts < MAX_CONTRAST_ATTEMPTS) {
        adjustedColor = useWhiteText
            ? adjustedColor.darken(0.2)
            : adjustedColor.brighten(0.2);
        attempts++;
    }
    return adjustedColor;
};

/**
 * 调整副标题颜色以确保可读性和和谐。
 * @param selectedColor 主选颜色。
 * @param useWhiteText 是否使用白色文本。
 * @param h 基础色的色相。
 * @param textRGB 文本的RGB颜色数组。
 * @param baseColor 基础颜色。
 * @returns 调整后的副标题颜色。
 */
const adjustSubtitleColor = (
    selectedColor: chroma.Color,
    useWhiteText: boolean,
    h: number,
    textRGB: number[],
    baseColor: chroma.Color
): chroma.Color => {
    let subtitleBase = useWhiteText
        ? selectedColor.brighten(0.9).desaturate(0.4).set("hsl.h", (h + SUBTITLE_HUE_SHIFT_WHITE) % HUE_MAX)
        : selectedColor.darken(0.5).desaturate(0.3).set("hsl.h", (h + SUBTITLE_HUE_SHIFT_BLACK) % HUE_MAX);

    let attempt = 0;
    while (attempt < MAX_SUBTITLE_CONTRAST_ATTEMPTS) {
        const contrastText = getContrast(subtitleBase.rgb(), textRGB);
        const contrastBg = getContrast(subtitleBase.rgb(), baseColor.rgb());
        if (contrastText >= MIN_SUBTITLE_CONTRAST_TEXT && contrastBg >= MIN_SUBTITLE_CONTRAST_BG) {
            break;
        }
        subtitleBase = useWhiteText
            ? subtitleBase.brighten(SUBTITLE_BRIGHTEN_DARKEN_AMOUNT).saturate(SUBTITLE_SATURATE_AMOUNT)
            : subtitleBase.darken(SUBTITLE_BRIGHTEN_DARKEN_AMOUNT).saturate(SUBTITLE_SATURATE_AMOUNT);
        attempt++;
    }
    return subtitleBase;
};

// ======================== 主题应用主逻辑 ======================== //

/**
 * 根据图片动态改变应用主题。
 * @param img HTMLImageElement 对象。
 * @returns 是否使用白色文本。
 */
export const changeTheme = (img: HTMLImageElement): boolean => {
    const baseColor = chroma(getColorSync(img).array());
    const [h, s, l] = baseColor.hsl();

    const useWhiteText = shouldUseWhiteText(img);
    const textRGB = useWhiteText ? [RGB_MAX_VALUE, RGB_MAX_VALUE, RGB_MAX_VALUE] : [0, 0, 0];
    const textColor = useWhiteText ? "#fff" : "#000";

    let selectedColor = adjustSelectedColor(baseColor, h, s, l);
    selectedColor = ensureContrastWithText(selectedColor, textRGB, useWhiteText);

    // 再次检查 selectedColor 和 baseColor 的对比度
    const contrastWithBase = getContrast(selectedColor.rgb(), baseColor.rgb());
    if (contrastWithBase < MIN_BASE_CONTRAST) {
        selectedColor = selectedColor.brighten(0.5);
    }

    // ========= 🎨 背景与辅助色 =========
    const backgroundBlendColor = chroma
        .mix(useWhiteText ? "#000" : "#fff", selectedColor, MIX_BLACK_WHITE_RATIO_BLEND)
        .set("hsl.l", `*${LUMINANCE_MULTIPLIER_BLEND}`)
        .desaturate(DESATURATE_AMOUNT_BLEND)
        .alpha(useWhiteText ? ALPHA_WHITE_TEXT_BLEND : ALPHA_BLACK_TEXT_BLEND)
        .css();

    const backgroundRightColor = chroma
        .mix(baseColor, selectedColor, MIX_BASE_SELECTED_RATIO_RIGHT)
        .set("hsl.s", `*${SATURATION_MULTIPLIER_RIGHT}`)
        .set("hsl.l", `*${LUMINANCE_MULTIPLIER_RIGHT}`)
        .alpha(useWhiteText ? ALPHA_WHITE_TEXT_RIGHT : ALPHA_BLACK_TEXT_RIGHT)
        .css();

    // ========= 副标题颜色更克制 =========
    const subtitleBase = adjustSubtitleColor(selectedColor, useWhiteText, h, textRGB, baseColor);
    const subtitleColor = subtitleBase.alpha(SUBTITLE_ALPHA).css();

    // -------- body-blur-color 更柔 --------
    const bodyBlurColor = chroma(useWhiteText ? "black" : "white")
        .alpha(useWhiteText ? BODY_BLUR_BLACK_ALPHA : BODY_BLUR_WHITE_ALPHA)
        .mix(baseColor, BODY_BLUR_MIX_RATIO)
        .desaturate(DESATURATE_AMOUNT_BLEND)
        .set("hsl.l", useWhiteText ? `*${BODY_BLUR_LUMINANCE_MULTIPLIER_WHITE}` : `*${BODY_BLUR_LUMINANCE_MULTIPLIER_BLACK}`)
        .css();

    // ========= ✅ 应用主题色 =========
    setCSSVariables({
        "text-color": textColor,
        "top-hr-color": subtitleColor,
        "left-item-selected-bg": selectedColor.css(),
        "blend-color": backgroundBlendColor,
        "right-bg-color": backgroundRightColor,
        "body-blur-color": bodyBlurColor,
    });

    return useWhiteText;
};

// ======================== 背景处理工具 ======================== //

/**
 * 从 CSS style 字符串中提取图片 URL。
 * @param style 包含 URL 的 CSS 字符串。
 * @returns 提取到的图片 URL，如果没有找到则为 null。
 */
const extractImageUrl = (style: string): string | null => {
    const match = style.match(/^url\(["']?(.*?)["']?\)$/);
    return match?.[1] || null;
};

let isBgLoading = false; // 标记背景是否正在加载中，避免重复请求

/**
 * 预加载背景图片并应用主题。
 * @param bg 背景图片 URL 或 CSS 渐变字符串。
 * @param cb 回调函数，在图片加载完成并应用主题后调用。
 */
export function preloadBackgroundImage(
    bg: string,
    cb: (bg: string, useWhite: boolean) => void
): void {
    if (isBgLoading) {
        console.warn("Background is loading, ignore new request:", bg);
        return;
    }

    // 如果 bg 不是一个图片 URL (例如是渐变色)，则直接回调
    if (!bg.startsWith("url(")) {
        cb(bg, false); // 假设非图片背景默认不使用白色文本
        return;
    }

    const imgUrl = extractImageUrl(bg);
    if (!imgUrl) {
        // 如果无法解析 URL，回退到默认背景
        return preloadBackgroundImage(DEFAULT_BACKGROUND_IMAGE, cb);
    }

    isBgLoading = true; // 设置加载中标记

    const img = new Image();
    let isResolved = false; // 标记是否已处理加载结果

    const timeoutId = setTimeout(() => {
        if (!isResolved) {
            console.error(`Background image load timed out: ${imgUrl}`);
            isResolved = true;
            isBgLoading = false;
            preloadBackgroundImage(DEFAULT_BACKGROUND_IMAGE, cb); // 超时回退到默认背景
        }
    }, IMAGE_LOAD_TIMEOUT);

    img.onload = () => {
        if (isResolved) return; // 避免重复处理
        clearTimeout(timeoutId);
        isResolved = true;
        isBgLoading = false;
        cb(bg, changeTheme(img)); // 图片加载成功，应用主题
    };

    img.onerror = () => {
        if (isResolved) return; // 避免重复处理
        clearTimeout(timeoutId);
        isResolved = true;
        console.error(`Failed to load background image: ${imgUrl}`);
        isBgLoading = false;
        preloadBackgroundImage(DEFAULT_BACKGROUND_IMAGE, cb); // 图片加载失败，回退到默认背景
    };

    img.src = imgUrl;
}