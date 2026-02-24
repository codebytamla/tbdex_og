import 'package:flutter/material.dart';

import '../theme/app_theme.dart';

/// Shared geometry for bottom chrome (nav bar + FAB) so their spacing stays
/// consistent when we tune sizes.
const kBottomChromeBottomInset = 20.0;
const kBottomChromeFabSize = 56.0;
const kBottomChromeFabRightInset = 20.0;
const kBottomChromeNavFabGap = 6.0;
const kBottomChromeNavRightInset =
    kBottomChromeFabRightInset + kBottomChromeFabSize + kBottomChromeNavFabGap;

/// Height of the nav bar pill.
const _kNavBarHeight = 56.0;

/// Horizontal padding on the items row inside the nav bar pill.
const _kNavBarHorizontalPadding = 5.0;

/// Vertical padding inside the nav bar pill — controls how tall the
/// inner selected-item highlight pill is. Smaller = taller active pill.
const _kNavBarVerticalPadding = 0.0;

const _kNavBarItemIconSize = 20.0;
const _kNavBarLabelSize = 9.0;

class LiquidGlassNavBar extends StatelessWidget {
  final int selectedIndex;
  final ValueChanged<int> onDestinationSelected;
  final List<LiquidGlassNavItem> items;

  const LiquidGlassNavBar({
    super.key,
    required this.selectedIndex,
    required this.onDestinationSelected,
    required this.items,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final navBarColor = isDark
        ? const Color(0xFF222326)
        : const Color(0xFFF6F7F8);

    return Positioned(
      left: 24,
      right: kBottomChromeNavRightInset,
      bottom: kBottomChromeBottomInset,
      child: Container(
        height: _kNavBarHeight,
        decoration: BoxDecoration(
          color: navBarColor,
          borderRadius: BorderRadius.circular(999),
          border: Border.all(
            color: isDark
                ? Colors.white.withValues(alpha: 0.14)
                : Colors.black.withValues(alpha: 0.08),
            width: 0.4,
          ),
        ),
        child: Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: _kNavBarHorizontalPadding,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: List.generate(items.length, (index) {
              return Expanded(
                child: _NavBarItem(
                  item: items[index],
                  isSelected: index == selectedIndex,
                  onTap: () => onDestinationSelected(index),
                ),
              );
            }),
          ),
        ),
      ),
    );
  }
}

class LiquidGlassNavItem {
  final IconData icon;
  final IconData selectedIcon;
  final String label;

  const LiquidGlassNavItem({
    required this.icon,
    required this.selectedIcon,
    required this.label,
  });
}

class _NavBarItem extends StatelessWidget {
  final LiquidGlassNavItem item;
  final bool isSelected;
  final VoidCallback onTap;

  const _NavBarItem({
    required this.item,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final primaryColor = isDark ? AppColorsDark.primary : AppColors.primary;

    final inactiveIconColor = isDark
        ? Colors.white.withValues(alpha: 0.78)
        : Colors.black.withValues(alpha: 0.74);
    final inactiveLabelColor = isDark
        ? Colors.white.withValues(alpha: 0.78)
        : Colors.black.withValues(alpha: 0.44);
    final iconColor = isSelected ? primaryColor : inactiveIconColor;
    final labelColor = isSelected ? primaryColor : inactiveLabelColor;

    final selectedBg = isDark
        ? Colors.white.withValues(alpha: 0.16)
        : Colors.black.withValues(alpha: 0.10);

    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: _kNavBarVerticalPadding),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 220),
          curve: Curves.easeOutCubic,
          decoration: BoxDecoration(
            color: isSelected ? selectedBg : Colors.transparent,
            borderRadius: BorderRadius.circular(999),
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                AnimatedSwitcher(
                  duration: const Duration(milliseconds: 180),
                  child: Icon(
                    isSelected ? item.selectedIcon : item.icon,
                    key: ValueKey(isSelected),
                    size: _kNavBarItemIconSize,
                    color: iconColor,
                  ),
                ),
                const SizedBox(height: 2),
                AnimatedDefaultTextStyle(
                  duration: const Duration(milliseconds: 180),
                  style: TextStyle(
                    fontSize: _kNavBarLabelSize,
                    fontWeight: isSelected ? FontWeight.w600 : FontWeight.w500,
                    color: labelColor,
                    letterSpacing: -0.1,
                  ),
                  child: Text(item.label, overflow: TextOverflow.ellipsis),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
