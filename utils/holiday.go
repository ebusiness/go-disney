package utils

import (
	"math"
	"time"
)

// IsHoliday - check the day is holiday or not
func IsHoliday(t time.Time) bool {
	return len((holiday{}).name(t)) > 0
}

type holiday struct{}

func (h holiday) name(t time.Time) string {
	t = t.In(jst)
	y, m, d, w := h.toYMDW(t)
	name := h.getName(y, m, d, w)

	if len(name) > 0 {
		return name
	}
	//振替休日
	if 1973 <= y && w == 0 {
		yname := h.getYesterdayNameFromTime(t)
		if len(yname) >= 1 {
			name = "振替休日"
		}
	} else if m == 5 && d == 6 && 2007 <= y && 1 <= w && w <= 2 {
		name = "振替休日"
	}
	return name
}

func (h holiday) toYMDW(t time.Time) (int, int, int, int) {
	return t.Year(), int(t.Month()), t.Day(), (int(t.Weekday()) + 6) % 7
}

func (h holiday) getYesterdayNameFromTime(t time.Time) string {
	yesterday := t.AddDate(0, 0, -1)
	y, m, d, w := h.toYMDW(yesterday)
	return h.getName(y, m, d, w)
}

func (h holiday) getName(y, m, d, w int) string {
	//皇室慶弔行事に伴う休日
	if y == 1959 && m == 4 && d == 10 {
		return "皇太子・明仁親王の結婚の儀"
	} else if y == 1989 && m == 2 && d == 24 {
		return "昭和天皇の大喪の礼"
	} else if y == 1990 && m == 11 && d == 12 {
		return "即位の礼正殿の儀"
	} else if y == 1993 && m == 6 && d == 9 {
		return "皇太子・徳仁親王の結婚の儀"
	} else if y == 2019 && m == 4 && d == 30 {
		return "国民の休日"
	} else if y == 2019 && m == 5 && d == 1 {
		return "天皇の即位の日"
	} else if y == 2019 && m == 5 && d == 2 {
		return "国民の休日"
	} else if y == 2019 && m == 10 && d == 22 {
		return "即位の礼正殿の儀"
	}
	olympic := h.getTokyoOlympic(y, m, d, w)
	if len(olympic) > 0 {
		return olympic
	}
	return h.getHolidayName(y, m, d, w)
}

func (h holiday) getTokyoOlympic(y, m, d, w int) string {
	if y != 2020 {
		return ""
	}
	// 東京オリンピック 特別措置法
	if m == 7 {
		if d == 23 {
			return "海の日"
		} else if d == 24 {
			return "スポーツの日"
		}
	} else if m == 8 {
		if d == 10 {
			return "山の日"
		}
	}
	return ""
}
func (h holiday) getHolidayName(y, m, d, w int) string {
	name := h.getHolidayNameOfJanuary(y, m, d, w)
	if len(name) == 0 {
		name = h.getHolidayNameOfFebruary(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfMarch(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfApril(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfMay(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfJune(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfJuly(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfAugust(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfSeptember(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfOctober(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfNovember(y, m, d, w)
	}
	if len(name) == 0 {
		name = h.getHolidayNameOfDecember(y, m, d, w)
	}
	return name
}
func (h holiday) getHolidayNameOfJanuary(y, m, d, w int) string {
	if m != 1 {
		return ""
	}
	if d == 1 {
		return "元日"
	}
	if 1949 <= y && y <= 1999 && d == 15 {
		return "成人の日"
	}

	if 2000 <= y && 8 <= d && d <= 14 && w == 0 {
		return "成人の日"
	}
	return ""
}
func (h holiday) getHolidayNameOfFebruary(y, m, d, w int) string {
	if m != 2 {
		return ""
	}

	if 2020 <= y && d == 23 {
		return "天皇誕生日"
	}
	if 1967 <= y && d == 11 {
		return "建国記念の日"
	}
	return ""
}
func (h holiday) getHolidayNameOfMarch(y, m, d, w int) string {
	if m != 3 {
		return ""
	}
	if 19 <= d && d <= 22 && d == h.shunBunDay(y) {
		return "春分の日"
	}
	return ""
}
func (h holiday) getHolidayNameOfApril(y, m, d, w int) string {
	if m != 4 || d != 29 {
		return ""
	}
	if y <= 1988 {
		return "天皇誕生日"
	} else if y <= 2006 {
		return "みどりの日"
	} else {
		return "昭和の日"
	}

}
func (h holiday) getHolidayNameOfMay(y, m, d, w int) string {
	if m != 5 {
		return ""
	}
	if d == 3 {
		return "憲法記念日"
	} else if d == 4 {
		if 1988 <= y && y <= 2006 && 1 <= w && w <= 5 {
			return "国民の休日"
		} else if 2007 <= y {
			return "みどりの日"
		}
	} else if d == 5 {
		return "こどもの日"
	}
	return ""
}

func (h holiday) getHolidayNameOfJune(y, m, d, w int) string {
	if m != 6 {
		return ""
	}
	return ""
}
func (h holiday) getHolidayNameOfJuly(y, m, d, w int) string {
	if m != 7 {
		return ""
	}
	if 1996 <= y && y <= 2002 {
		if d == 20 {
			return "海の日"
		}
	} else if 2003 <= y {
		if 15 <= d && d <= 21 && w == 0 && y != 2020 {
			return "海の日"
		}
	}
	return ""
}
func (h holiday) getHolidayNameOfAugust(y, m, d, w int) string {
	if m != 8 {
		return ""
	}
	if 2016 <= y {
		if d == 11 && y != 2020 {
			return "山の日"
		}
	}
	return ""
}
func (h holiday) getHolidayNameOfSeptember(y, m, d, w int) string {
	if m != 9 {
		return ""
	}
	name := h.getHolidayNameOfSeptember1(y, m, d, w)
	if len(name) == 0 {
		name = h.getHolidayNameOfSeptember2(y, m, d, w)
	}
	return name
}
func (h holiday) getHolidayNameOfSeptember1(y, m, d, w int) string {
	if m != 9 {
		return ""
	}
	if 1966 <= y && y <= 2002 {
		if d == 15 {
			return "敬老の日"
		}
	} else if 2003 <= y {
		if 15 <= d && d <= 21 && w == 0 {
			return "敬老の日"
		}
	}
	return ""
}
func (h holiday) getHolidayNameOfSeptember2(y, m, d, w int) string {
	if m != 9 {
		return ""
	}
	if 2009 <= y && w == 1 {
		if 21 <= d && d <= 23 {
			if d+1 == h.shuuBunDay(y) {
				return "国民の休日"
			}
		}
	}
	if 22 <= d && d <= 24 {
		if d == h.shuuBunDay(y) {
			return "秋分の日"
		}
	}
	return ""
}
func (h holiday) getHolidayNameOfOctober(y, m, d, w int) string {
	if m != 10 {
		return ""
	}
	if 1966 <= y && y <= 1999 {
		if d == 10 {
			return "体育の日"
		}
	} else if 2000 <= y {
		if 8 <= d && d <= 14 && w == 0 && y != 2020 {
			if 2020 <= y {
				return "スポーツの日"
			}
			return "体育の日"
		}
	}
	return ""
}
func (h holiday) getHolidayNameOfNovember(y, m, d, w int) string {
	if m != 11 {
		return ""
	}
	if d == 3 {
		return "文化の日"
	} else if d == 23 {
		return "勤労感謝の日"
	}
	return ""
}
func (h holiday) getHolidayNameOfDecember(y, m, d, w int) string {
	if m != 12 {
		return ""
	}
	if 1989 <= y && y <= 2018 && d == 23 {
		return "天皇誕生日"
	}
	return ""
}

func (h holiday) calcDay(future, present, past float64, y int) int {
	add := 0.242194*float64(y-1980) - math.Floor(float64(y-1980)/4.0)
	val := 0.0

	switch {
	case 2100 <= y && y <= 2150:
		val = future + add
	case 1980 <= y:
		val = present + add
	case 1900 <= y:
		val = past + add
	}

	return int(math.Floor(val))
}

func (h holiday) shunBunDay(y int) int {
	return h.calcDay(21.8519, 20.8431, 20.8357, y)
}

func (h holiday) shuuBunDay(y int) int {
	return h.calcDay(24.2488, 23.2488, 23.2588, y)
}
