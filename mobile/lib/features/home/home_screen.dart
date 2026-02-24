import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';
import '../../core/models/models.dart';
import '../../../main.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final tasks = MockData.todaysTasks;
    const previewTaskCount = 3;
    final previewTasks = tasks.take(previewTaskCount).toList();
    final hasMoreTasks = tasks.length > previewTaskCount;
    final assessmentCount = tasks
        .where((t) => t.status == 'pending' || t.status == 'in_progress')
        .length;
    final referralCount = 3;
    final totalPatients = MockData.patients.length;
    final bottomInset = MediaQuery.of(context).padding.bottom;
    final listBottomSpacing = bottomInset + 84;
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: isDark ? AppColorsDark.background : AppColors.background,
      body: SafeArea(
        bottom: false,
        child: CustomScrollView(
          slivers: [
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(22, 16, 22, 0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // ── Header ──────────────────────────────────────────────
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              // Text(
                              //   'Today',
                              //   style:
                              //       Theme.of(context).textTheme.bodyMedium?.copyWith(
                              //             fontWeight: FontWeight.w500,
                              //           ),
                              // ),
                              const SizedBox(height: 2),
                              Text(
                                'Overview',
                                style: Theme.of(context).textTheme.displaySmall
                                    ?.copyWith(
                                      fontWeight: FontWeight.w800,
                                      letterSpacing: -0.5,
                                      height: 1.1,
                                    ),
                              ),
                            ],
                          ),
                        ),
                        ValueListenableBuilder<ThemeMode>(
                          valueListenable: themeNotifier,
                          builder: (context, mode, child) => IconButton(
                            onPressed: () {
                              themeNotifier.value = mode == ThemeMode.light
                                  ? ThemeMode.dark
                                  : ThemeMode.light;
                            },
                            icon: Icon(
                              mode == ThemeMode.light
                                  ? Icons.dark_mode_outlined
                                  : Icons.light_mode_outlined,
                              size: 22,
                              color: isDark
                                  ? AppColorsDark.textSecondary
                                  : AppColors.textSecondary,
                            ),
                          ),
                        ),
                        const SizedBox(width: 2),
                        // Amber avatar
                        Container(
                          width: 40,
                          height: 40,
                          decoration: BoxDecoration(
                            color: const Color(0xFFD97706),
                            borderRadius: BorderRadius.circular(999),
                          ),
                          alignment: Alignment.center,
                          child: const Text(
                            'TD',
                            style: TextStyle(
                              color: Colors.white,
                              fontWeight: FontWeight.w700,
                              fontSize: 13,
                            ),
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 18),

                    // ── Sync Banner ─────────────────────────────────────────
                    Material(
                      color: Colors.transparent,
                      child: InkWell(
                        onTap: () => context.push('/sync'),
                        borderRadius: BorderRadius.circular(14),
                        child: Ink(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 14,
                            vertical: 7,
                          ),
                          decoration: BoxDecoration(
                            color: isDark
                                ? AppColorsDark.warningLight.withValues(
                                    alpha: 0.22,
                                  )
                                : const Color(0xFFFFF8EE),
                            borderRadius: BorderRadius.circular(14),
                            border: Border.all(
                              color: isDark
                                  ? const Color(
                                      0xFF92400E,
                                    ).withValues(alpha: 0.45)
                                  : const Color(0xFFFFE4B5),
                            ),
                          ),
                          child: Row(
                            children: [
                              Container(
                                width: 28,
                                height: 28,
                                decoration: BoxDecoration(
                                  color: isDark
                                      ? const Color(
                                          0xFFD97706,
                                        ).withValues(alpha: 0.18)
                                      : const Color(0xFFFFF3CD),
                                  borderRadius: BorderRadius.circular(10),
                                ),
                                child: Icon(
                                  Icons.sync_rounded,
                                  color: isDark
                                      ? const Color(0xFFFBBF24)
                                      : const Color(0xFFD97706),
                                  size: 18,
                                ),
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Text(
                                      'Tap to sync data',
                                      style: TextStyle(
                                        color: isDark
                                            ? const Color(0xFFD4B97A)
                                            : const Color(0xFF92400E),
                                        fontWeight: FontWeight.w600,
                                        fontSize: 13,
                                      ),
                                    ),
                                    const SizedBox(height: 1),
                                    Text(
                                      'Last synced 2 hours ago',
                                      style: TextStyle(
                                        color: isDark
                                            ? const Color(
                                                0xFFD4B97A,
                                              ).withValues(alpha: 0.55)
                                            : const Color(
                                                0xFF92400E,
                                              ).withValues(alpha: 0.65),
                                        fontWeight: FontWeight.w400,
                                        fontSize: 11.5,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              Icon(
                                Icons.chevron_right_rounded,
                                color: isDark
                                    ? const Color(0xFFC4A060)
                                    : const Color(0xFFD97706),
                                size: 20,
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 14),

                    // ── Hero: Total Patients ─────────────────────────────────
                    _HeroPatientCard(totalPatients: totalPatients),

                    const SizedBox(height: 12),

                    // ── Stats Row ────────────────────────────────────────────
                    Row(
                      children: [
                        Expanded(
                          child: _MiniStatCard(
                            label: 'ASSESSMENTS',
                            value: '$assessmentCount',
                            icon: Icons.description_outlined,
                            iconColor: const Color(0xFFE05C5C),
                            iconColorDark: const Color(0xFFA07070),
                            iconBg: const Color(0xFFFEECEC),
                            iconBgDark: const Color(
                              0xFF5C2A2A,
                            ).withValues(alpha: 0.3),
                            isDark: isDark,
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: _MiniStatCard(
                            label: 'REFERRALS',
                            value: '$referralCount',
                            icon: Icons.person_add_outlined,
                            iconColor: const Color(0xFF3AAFA9),
                            iconColorDark: const Color(0xFF5E9E9A),
                            iconBg: const Color(0xFFE6F7F7),
                            iconBgDark: const Color(
                              0xFF1A3A38,
                            ).withValues(alpha: 0.4),
                            isDark: isDark,
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 36),

                    // ── Recent Assessments header ────────────────────────────
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Text(
                          'Recents',
                          style: Theme.of(context).textTheme.headlineSmall
                              ?.copyWith(fontWeight: FontWeight.w700),
                        ),
                        if (hasMoreTasks)
                          InkWell(
                            onTap: () => context.go('/activity'),
                            borderRadius: BorderRadius.circular(8),
                            child: Padding(
                              padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Text(
                                    'See All',
                                    style: TextStyle(
                                      color: Theme.of(context).colorScheme.primary,
                                      fontWeight: FontWeight.w600,
                                      fontSize: 14,
                                    ),
                                  ),
                                  Icon(
                                    Icons.chevron_right_rounded,
                                    size: 18,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ],
                              ),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 12),
                  ],
                ),
              ),
            ),

            // ── Task list ────────────────────────────────────────────────────
            SliverPadding(
              padding: const EdgeInsets.symmetric(horizontal: 22),
              sliver: SliverList.separated(
                itemCount: previewTasks.length,
                separatorBuilder: (_, _) => const SizedBox(height: 10),
                itemBuilder: (context, index) {
                  final task = previewTasks[index];
                  return _TaskRow(
                    task: task,
                    onTap: () => context.push('/patients/${task.patientId}'),
                  );
                },
              ),
            ),

            SliverToBoxAdapter(child: SizedBox(height: listBottomSpacing)),
          ],
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Hero Patient Card
// ─────────────────────────────────────────────────────────────────────────────
class _HeroPatientCard extends StatelessWidget {
  final int totalPatients;
  const _HeroPatientCard({required this.totalPatients});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final cardColor = isDark ? AppColorsDark.surface : AppColors.surface;
    final borderColor = isDark ? AppColorsDark.border : AppColors.border;
    final textPrimary = isDark
        ? AppColorsDark.textPrimary
        : AppColors.textPrimary;
    final textSecondary = isDark
        ? AppColorsDark.textSecondary
        : AppColors.textSecondary;
    final iconBg = isDark
        ? AppColorsDark.primaryLight.withValues(alpha: 0.35)
        : AppColors.primaryLight;

    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: cardColor,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: borderColor),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: isDark ? 0.18 : 0.04),
            blurRadius: 10,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Person-group icon in a subtle box
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(Icons.group_outlined, color: Colors.white, size: 22),
          ),

          const SizedBox(height: 20),

          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              // Number + label
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '$totalPatients',
                      style: TextStyle(
                        color: textPrimary,
                        fontSize: 42,
                        fontWeight: FontWeight.w700,
                        height: 1,
                      ),
                    ),
                    const SizedBox(height: 6),
                    Text(
                      'TOTAL PATIENTS',
                      style: TextStyle(
                        color: textSecondary,
                        fontSize: 11,
                        fontWeight: FontWeight.w600,
                        letterSpacing: 1.1,
                      ),
                    ),
                  ],
                ),
              ),

              // Mini bar chart decoration
              const _BarChartDecoration(),
            ],
          ),
        ],
      ),
    );
  }
}

// Small decorative bar chart
class _BarChartDecoration extends StatelessWidget {
  const _BarChartDecoration();

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final bars = isDark
        ? const [
            _Bar(height: 28, color: Color(0xFFBFA8A8)),
            _Bar(height: 20, color: Color(0xFF9BBCBA)),
            _Bar(height: 38, color: Color(0xFFC4A97D)),
            _Bar(height: 24, color: Color(0xFF9BBCBA)),
          ]
        : const [
            _Bar(height: 28, color: Color(0xFFE05C5C)),
            _Bar(height: 20, color: Color(0xFF3AAFA9)),
            _Bar(height: 38, color: Color(0xFFD97706)),
            _Bar(height: 24, color: Color(0xFF3AAFA9)),
          ];

    return Row(
      crossAxisAlignment: CrossAxisAlignment.end,
      children: bars
          .map(
            (b) => Padding(
              padding: const EdgeInsets.only(left: 5),
              child: Container(
                width: 10,
                height: b.height,
                decoration: BoxDecoration(
                  color: b.color.withValues(alpha: isDark ? 0.7 : 0.85),
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
            ),
          )
          .toList(),
    );
  }
}

class _Bar {
  final double height;
  final Color color;
  const _Bar({required this.height, required this.color});
}

// ─────────────────────────────────────────────────────────────────────────────
// Mini Stat Card (Assessments / Referrals)
// ─────────────────────────────────────────────────────────────────────────────
class _MiniStatCard extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color iconColor;
  final Color iconColorDark;
  final Color iconBg;
  final Color iconBgDark;
  final bool isDark;

  const _MiniStatCard({
    required this.label,
    required this.value,
    required this.icon,
    required this.iconColor,
    required this.iconColorDark,
    required this.iconBg,
    required this.iconBgDark,
    required this.isDark,
  });

  @override
  Widget build(BuildContext context) {
    final effectiveColor = isDark ? iconColorDark : iconColor;
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isDark ? AppColorsDark.surface : AppColors.surface,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(
          color: isDark ? AppColorsDark.border : AppColors.border,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: isDark ? 0.18 : 0.04),
            blurRadius: 10,
            offset: const Offset(0, 3),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: isDark ? iconBgDark : iconBg,
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: effectiveColor, size: 18),
          ),
          const SizedBox(height: 16),
          Text(
            value,
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.w700,
              color: effectiveColor,
              height: 1,
            ),
          ),
          const SizedBox(height: 5),
          Text(
            label,
            style: TextStyle(
              fontSize: 10.5,
              fontWeight: FontWeight.w600,
              color: isDark
                  ? AppColorsDark.textSecondary
                  : AppColors.textSecondary,
              letterSpacing: 0.9,
            ),
          ),
        ],
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Task Row
// ─────────────────────────────────────────────────────────────────────────────
class _TaskRow extends StatelessWidget {
  final TaskItem task;
  final VoidCallback onTap;

  const _TaskRow({required this.task, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        borderRadius: BorderRadius.circular(20),
        onTap: onTap,
        child: Ink(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
          decoration: BoxDecoration(
            color: isDark ? AppColorsDark.surface : AppColors.surface,
            borderRadius: BorderRadius.circular(20),
            border: Border.all(
              color: isDark ? AppColorsDark.border : AppColors.border,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: isDark ? 0.14 : 0.03),
                blurRadius: 8,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      task.patientName,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const SizedBox(height: 3),
                    Text(
                      '${task.village} • ${task.time}',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.w500,
                        color: isDark
                            ? AppColorsDark.textSecondary
                            : AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 10),
              _StatusChip(status: task.status),
              const SizedBox(width: 6),
              Icon(
                Icons.chevron_right_rounded,
                size: 18,
                color: isDark
                    ? AppColorsDark.textTertiary
                    : AppColors.textTertiary,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────────────────────────────────────
// Status Chip
// ─────────────────────────────────────────────────────────────────────────────
class _StatusChip extends StatelessWidget {
  final String status;
  const _StatusChip({required this.status});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final style = _StatusStyle.resolve(status, isDark);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      decoration: BoxDecoration(
        color: style.background,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        style.label,
        style: TextStyle(
          color: style.foreground,
          fontSize: 11,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.3,
        ),
      ),
    );
  }
}

class _StatusStyle {
  final Color foreground;
  final Color background;
  final String label;

  const _StatusStyle({
    required this.foreground,
    required this.background,
    required this.label,
  });

  static _StatusStyle resolve(String status, bool isDark) {
    switch (status) {
      case 'pending':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFFDE68A)
              : const Color(0xFF92400E),
          background: isDark
              ? AppColorsDark.warningLight.withValues(alpha: 0.35)
              : AppColors.warningLight,
          label: 'PENDING',
        );
      case 'in_progress':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFBFDBFE)
              : const Color(0xFF1E40AF),
          background: isDark
              ? AppColorsDark.primaryLight.withValues(alpha: 0.38)
              : AppColors.primaryLight,
          label: 'IN PROGRESS',
        );
      case 'done':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFF6EE7B7)
              : const Color(0xFF065F46),
          background: isDark
              ? AppColorsDark.successLight.withValues(alpha: 0.35)
              : AppColors.successLight,
          label: 'DONE',
        );
      default:
        return _StatusStyle(
          foreground: isDark
              ? AppColorsDark.textSecondary
              : AppColors.textSecondary,
          background: isDark ? AppColorsDark.surfaceDim : AppColors.surfaceDim,
          label: status.toUpperCase(),
        );
    }
  }
}
