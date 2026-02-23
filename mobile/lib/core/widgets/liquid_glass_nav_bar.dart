import 'dart:ui';

import 'package:flutter/material.dart';

import '../theme/app_theme.dart';

/// Height of the nav bar pill.
const _kNavBarHeight = 62.0;

/// Horizontal padding on the items row inside the nav bar pill.
const _kNavBarHorizontalPadding = 6.0;

/// Vertical padding inside the nav bar pill — controls how tall the
/// inner selected-item highlight pill is. Smaller = taller active pill.
const _kNavBarVerticalPadding = 7.0;

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

    // Glass gradient — lighter at top (catches light), darker at bottom.
    final pillGradient = LinearGradient(
      begin: Alignment.topCenter,
      end: Alignment.bottomCenter,
      colors: isDark
          ? [
              const Color(0xFF3C3C3E).withValues(alpha: 0.45),
              const Color(0xFF1C1C1E).withValues(alpha: 0.55),
            ]
          : [
              Colors.white.withValues(alpha: 0.55),
              Colors.white.withValues(alpha: 0.38),
            ],
    );

    return Positioned(
      left: 24,
      right: 88,
      bottom: 20,
      child: ClipRRect(
        borderRadius: BorderRadius.circular(999),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 40, sigmaY: 40),
          child: Container(
            height: _kNavBarHeight,
            decoration: BoxDecoration(
              gradient: pillGradient,
              borderRadius: BorderRadius.circular(999),
              border: Border.all(
                color: isDark
                    ? Colors.white.withValues(alpha: 0.12)
                    : Colors.white.withValues(alpha: 0.70),
                width: 0.8,
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.30),
                  blurRadius: 30,
                  offset: const Offset(0, 10),
                ),
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.10),
                  blurRadius: 6,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: Stack(
              children: [
                // Specular highlight — simulates light refracting off the
                // top surface of the glass pill.
                Positioned(
                  left: 0,
                  right: 0,
                  top: 0,
                  height: 22,
                  child: DecoratedBox(
                    decoration: BoxDecoration(
                      borderRadius: const BorderRadius.vertical(
                        top: Radius.circular(999),
                      ),
                      gradient: LinearGradient(
                        begin: Alignment.topCenter,
                        end: Alignment.bottomCenter,
                        colors: [
                          Colors.white.withValues(alpha: isDark ? 0.10 : 0.55),
                          Colors.white.withValues(alpha: 0),
                        ],
                      ),
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(
                      horizontal: _kNavBarHorizontalPadding),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: List.generate(items.length, (index) {
                      return Expanded(
                        child: _NavBarItem(
                          item: items[index],
                          isSelected: index == selectedIndex,
                          isDark: isDark,
                          onTap: () => onDestinationSelected(index),
                        ),
                      );
                    }),
                  ),
                ),
              ],
            ),
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
  final bool isDark;
  final VoidCallback onTap;

  const _NavBarItem({
    required this.item,
    required this.isSelected,
    required this.isDark,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final primaryColor = isDark ? AppColorsDark.primary : AppColors.primary;

    final inactiveColor = isDark
        ? Colors.white.withValues(alpha: 0.75)
        : Colors.black.withValues(alpha: 0.28);
    final iconColor = isSelected ? primaryColor : inactiveColor;
    final labelColor = isSelected ? primaryColor : inactiveColor;

    final selectedBg = isDark
        ? Colors.white.withValues(alpha: 0.10)
        : Colors.black.withValues(alpha: 0.06);

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
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                AnimatedSwitcher(
                  duration: const Duration(milliseconds: 180),
                  child: Icon(
                    isSelected ? item.selectedIcon : item.icon,
                    key: ValueKey(isSelected),
                    size: 22,
                    color: iconColor,
                  ),
                ),
                const SizedBox(height: 2),
                AnimatedDefaultTextStyle(
                  duration: const Duration(milliseconds: 180),
                  style: TextStyle(
                    fontSize: 10,
                    fontWeight:
                        isSelected ? FontWeight.w600 : FontWeight.w500,
                    color: labelColor,
                    letterSpacing: -0.1,
                  ),
                  child: Text(
                    item.label,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
