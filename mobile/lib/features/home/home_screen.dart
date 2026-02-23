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
    const previewTaskCount = 4;
    final previewTasks = tasks.take(previewTaskCount).toList();
    final hasMoreTasks = tasks.length > previewTaskCount;
    final pendingCount = tasks.where((t) => t.status == 'pending').length;
    final bottomInset = MediaQuery.of(context).padding.bottom;
    final listBottomSpacing = bottomInset + 128;
    final isDark = Theme.of(context).brightness == Brightness.dark;
    Color lightStatBackground(Color accent, double alpha) =>
        Color.alphaBlend(accent.withValues(alpha: alpha), AppColors.surface);
    Color lightStatBorder(Color accent, double alpha) =>
        Color.alphaBlend(accent.withValues(alpha: alpha), AppColors.border);
    Color darkStatBackground(Color accent, double alpha) => Color.alphaBlend(
      accent.withValues(alpha: alpha),
      AppColorsDark.surface,
    );
    Color darkStatBorder(Color accent, double alpha) =>
        Color.alphaBlend(accent.withValues(alpha: alpha), AppColorsDark.border);

    return Scaffold(
      backgroundColor: isDark ? AppColorsDark.background : AppColors.background,
      body: SafeArea(
        child: CustomScrollView(
          slivers: [
            SliverToBoxAdapter(
              child: Padding(
                padding: const EdgeInsets.fromLTRB(24, 16, 24, 0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Good morning,',
                                style: Theme.of(context).textTheme.bodyMedium
                                    ?.copyWith(fontWeight: FontWeight.w500),
                              ),
                              const SizedBox(height: 2),
                              Text(
                                'Sarah Kimani',
                                style: Theme.of(context)
                                    .textTheme
                                    .headlineMedium
                                    ?.copyWith(
                                      fontWeight: FontWeight.w700,
                                      letterSpacing: -0.3,
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
                                  ? Icons.dark_mode
                                  : Icons.light_mode,
                              color: isDark
                                  ? AppColorsDark.textSecondary
                                  : AppColors.textSecondary,
                            ),
                            tooltip: mode == ThemeMode.light
                                ? 'Switch to dark mode'
                                : 'Switch to light mode',
                          ),
                        ),
                        const SizedBox(width: 2),
                        Container(
                          width: 40,
                          height: 40,
                          decoration: BoxDecoration(
                            color: isDark
                                ? AppColorsDark.primaryLight.withValues(
                                    alpha: 0.35,
                                  )
                                : AppColors.primaryLight,
                            borderRadius: BorderRadius.circular(999),
                            border: Border.all(
                              color: isDark
                                  ? AppColors.primary.withValues(alpha: 0.5)
                                  : const Color(0xFFBFDBFE),
                            ),
                          ),
                          alignment: Alignment.center,
                          child: const Text(
                            'SK',
                            style: TextStyle(
                              color: AppColors.primary,
                              fontWeight: FontWeight.w700,
                              fontSize: 13,
                            ),
                          ),
                        ),
                      ],
                    ),

                    const SizedBox(height: 16),

                    Material(
                      color: Colors.transparent,
                      child: InkWell(
                        onTap: () => context.push('/sync'),
                        borderRadius: BorderRadius.circular(12),
                        child: Ink(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 14,
                            vertical: 12,
                          ),
                          decoration: BoxDecoration(
                            color: isDark
                                ? AppColorsDark.warningLight.withValues(
                                    alpha: 0.28,
                                  )
                                : const Color(0xFFFFFBEB),
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: isDark
                                  ? const Color(
                                      0xFF92400E,
                                    ).withValues(alpha: 0.5)
                                  : const Color(0xFFFDE68A),
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: isDark
                                    ? Colors.black.withValues(alpha: 0.2)
                                    : const Color(
                                        0xFFF59E0B,
                                      ).withValues(alpha: 0.06),
                                blurRadius: 12,
                                offset: const Offset(0, 4),
                              ),
                            ],
                          ),
                          child: Row(
                            children: [
                              Icon(
                                Icons.sync_rounded,
                                color: isDark
                                    ? const Color(0xFFFBBF24)
                                    : const Color(0xFF92400E),
                                size: 18,
                              ),
                              const SizedBox(width: 10),
                              Expanded(
                                child: Text(
                                  'Last synced: 2 hours ago',
                                  style: TextStyle(
                                    color: isDark
                                        ? const Color(0xFFFDE68A)
                                        : const Color(0xFF92400E),
                                    fontWeight: FontWeight.w600,
                                    fontSize: 13,
                                  ),
                                ),
                              ),
                              Icon(
                                Icons.chevron_right_rounded,
                                color: isDark
                                    ? const Color(0xFFF59E0B)
                                    : const Color(0xFFD97706),
                                size: 20,
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 20),

                    Row(
                      children: [
                        _StatCard(
                          value: '$pendingCount',
                          label: 'Pending\nAssessments',
                          color: AppColors.primary,
                          bgColor: isDark
                              ? darkStatBackground(AppColors.primary, 0.14)
                              : lightStatBackground(AppColors.primary, 0.1),
                          borderColor: isDark
                              ? darkStatBorder(AppColors.primary, 0.3)
                              : lightStatBorder(AppColors.primary, 0.22),
                        ),
                        const SizedBox(width: 10),
                        _StatCard(
                          value: '3',
                          label: 'Referrals\nToday',
                          color: AppColors.danger,
                          bgColor: isDark
                              ? darkStatBackground(AppColors.danger, 0.13)
                              : lightStatBackground(AppColors.danger, 0.09),
                          borderColor: isDark
                              ? darkStatBorder(AppColors.danger, 0.28)
                              : lightStatBorder(AppColors.danger, 0.2),
                        ),
                        const SizedBox(width: 10),
                        _StatCard(
                          value: '${MockData.patients.length}',
                          label: 'Total\nPatients',
                          color: AppColors.success,
                          bgColor: isDark
                              ? darkStatBackground(AppColors.success, 0.12)
                              : lightStatBackground(AppColors.success, 0.08),
                          borderColor: isDark
                              ? darkStatBorder(AppColors.success, 0.28)
                              : lightStatBorder(AppColors.success, 0.2),
                        ),
                      ],
                    ),

                    const SizedBox(height: 24),

                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Text(
                          "Today's Tasks",
                          style: Theme.of(context).textTheme.headlineSmall
                              ?.copyWith(fontWeight: FontWeight.w700),
                        ),
                        if (hasMoreTasks)
                          TextButton(
                            onPressed: () => context.go('/activity'),
                            style: TextButton.styleFrom(
                              padding: EdgeInsets.zero,
                              minimumSize: const Size(52, 30),
                              tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                              foregroundColor: isDark
                                  ? AppColorsDark.primary
                                  : AppColors.primary,
                            ),
                            child: const Text(
                              'See All',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 14,
                              ),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 14),
                  ],
                ),
              ),
            ),

            SliverPadding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              sliver: SliverList.separated(
                itemCount: previewTasks.length,
                separatorBuilder: (_, _) => const SizedBox(height: 12),
                itemBuilder: (context, index) {
                  final task = previewTasks[index];
                  return _TaskCard(
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

class _StatCard extends StatelessWidget {
  final String value;
  final String label;
  final Color color;
  final Color bgColor;
  final Color borderColor;

  const _StatCard({
    required this.value,
    required this.label,
    required this.color,
    required this.bgColor,
    required this.borderColor,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: borderColor),
          boxShadow: [
            BoxShadow(
              color: isDark
                  ? Colors.black.withValues(alpha: 0.25)
                  : Colors.black.withValues(alpha: 0.04),
              blurRadius: 12,
              offset: const Offset(0, 4),
            ),
          ],
        ),
        child: Column(
          children: [
            Text(
              value,
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.w700,
                color: color,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              label,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 11,
                fontWeight: FontWeight.w500,
                color: color.withValues(alpha: 0.8),
                height: 1.3,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TaskCard extends StatelessWidget {
  final TaskItem task;
  final VoidCallback onTap;

  const _TaskCard({required this.task, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Material(
      color: Colors.transparent,
      child: InkWell(
        borderRadius: BorderRadius.circular(16),
        onTap: onTap,
        child: Ink(
          padding: const EdgeInsets.all(15),
          decoration: BoxDecoration(
            color: isDark ? AppColorsDark.surface : AppColors.surface,
            borderRadius: BorderRadius.circular(16),
            border: Border.all(
              color: isDark ? AppColorsDark.border : AppColors.border,
            ),
            boxShadow: [
              BoxShadow(
                color: isDark
                    ? Colors.black.withValues(alpha: 0.16)
                    : Colors.black.withValues(alpha: 0.04),
                blurRadius: 10,
                offset: const Offset(0, 3),
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
                    const SizedBox(height: 4),
                    Text(
                      '${task.village} • ${task.time}',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                ),
              ),
              _StatusChip(status: task.status),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatusChip extends StatelessWidget {
  final String status;

  const _StatusChip({required this.status});

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final style = _StatusStyle.resolve(status, isDark);

    return AnimatedContainer(
      duration: const Duration(milliseconds: 220),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: style.background,
        border: Border.all(color: style.border),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        style.label,
        style: TextStyle(
          color: style.foreground,
          fontSize: 12,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class _StatusStyle {
  final Color foreground;
  final Color background;
  final Color border;
  final String label;

  const _StatusStyle({
    required this.foreground,
    required this.background,
    required this.border,
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
          border: isDark
              ? AppColors.warning.withValues(alpha: 0.45)
              : const Color(0xFFFDE68A),
          label: 'Pending',
        );
      case 'in_progress':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFFBFDBFE)
              : const Color(0xFF1E40AF),
          background: isDark
              ? AppColorsDark.primaryLight.withValues(alpha: 0.38)
              : AppColors.primaryLight,
          border: isDark
              ? AppColors.primary.withValues(alpha: 0.45)
              : const Color(0xFFBFDBFE),
          label: 'In Progress',
        );
      case 'done':
        return _StatusStyle(
          foreground: isDark
              ? const Color(0xFF6EE7B7)
              : const Color(0xFF065F46),
          background: isDark
              ? AppColorsDark.successLight.withValues(alpha: 0.35)
              : AppColors.successLight,
          border: isDark
              ? AppColors.success.withValues(alpha: 0.4)
              : const Color(0xFFA7F3D0),
          label: 'Done',
        );
      default:
        return _StatusStyle(
          foreground: isDark
              ? AppColorsDark.textSecondary
              : AppColors.textSecondary,
          background: isDark ? AppColorsDark.surfaceDim : AppColors.surfaceDim,
          border: isDark ? AppColorsDark.border : AppColors.border,
          label: status,
        );
    }
  }
}
