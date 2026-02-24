import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';

class ActivityScreen extends StatefulWidget {
  const ActivityScreen({super.key});

  @override
  State<ActivityScreen> createState() => _ActivityScreenState();
}

class _ActivityScreenState extends State<ActivityScreen> {
  String _filter = 'all';

  @override
  Widget build(BuildContext context) {
    var assessments = MockData.assessmentHistory;
    if (_filter != 'all') {
      assessments = assessments.where((a) => a.riskLevel == _filter).toList();
    }

    return Scaffold(
      backgroundColor: context.appBackground,
      appBar: AppBar(
        title: const Text('Activity'),
        actions: [
          IconButton(
            icon: const Icon(Icons.sync_rounded),
            tooltip: 'Sync Center',
            onPressed: () => context.push('/sync'),
          ),
        ],
      ),
      body: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Sync summary bar
          Container(
            margin: const EdgeInsets.fromLTRB(20, 8, 20, 0),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: context.appPrimaryLight,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Row(
              children: [
                Icon(
                  Icons.cloud_upload_rounded,
                  size: 18,
                  color: context.appPrimary,
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    '8 assessments pending sync',
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      color: context.appPrimary,
                    ),
                  ),
                ),
                InkWell(
                  onTap: () => context.push('/sync'),
                  borderRadius: BorderRadius.circular(4),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
                    child: Text(
                      'View',
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: context.appPrimary,
                        decoration: TextDecoration.underline,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),

          const SizedBox(height: 16),

          // Filter chips
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Row(
              children: [
                _FilterChip(
                  label: 'All',
                  selected: _filter == 'all',
                  onTap: () => setState(() => _filter = 'all'),
                ),
                const SizedBox(width: 8),
                _FilterChip(
                  label: 'High Risk',
                  selected: _filter == 'high',
                  onTap: () => setState(() => _filter = 'high'),
                  color: AppColors.danger,
                ),
                const SizedBox(width: 8),
                _FilterChip(
                  label: 'Medium',
                  selected: _filter == 'medium',
                  onTap: () => setState(() => _filter = 'medium'),
                  color: AppColors.warning,
                ),
                const SizedBox(width: 8),
                _FilterChip(
                  label: 'Low',
                  selected: _filter == 'low',
                  onTap: () => setState(() => _filter = 'low'),
                  color: AppColors.success,
                ),
              ],
            ),
          ),

          const SizedBox(height: 8),

          // Section header
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 8, 20, 4),
            child: Text(
              'Recent Assessments',
              style: Theme.of(context).textTheme.titleLarge,
            ),
          ),

          // Assessment list
          Expanded(
            child: assessments.isEmpty
                ? Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          Icons.assessment_outlined,
                          size: 48,
                          color: context.appTextTertiary,
                        ),
                        const SizedBox(height: 12),
                        Text(
                          'No assessments found',
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ],
                    ),
                  )
                : ListView.separated(
                    padding: const EdgeInsets.fromLTRB(20, 8, 20, 100),
                    itemCount: assessments.length,
                    separatorBuilder: (_, __) => const SizedBox(height: 8),
                    itemBuilder: (context, index) {
                      final a = assessments[index];
                      return _AssessmentHistoryCard(assessment: a);
                    },
                  ),
          ),
        ],
      ),
    );
  }
}

class _FilterChip extends StatelessWidget {
  final String label;
  final bool selected;
  final VoidCallback onTap;
  final Color? color;

  const _FilterChip({
    required this.label,
    required this.selected,
    required this.onTap,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    final chipColor = color ?? context.appPrimary;

    return Material(
      color: Colors.transparent,
      borderRadius: BorderRadius.circular(20),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(20),
        child: Ink(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
          decoration: BoxDecoration(
            color: selected ? chipColor : context.appSurface,
            borderRadius: BorderRadius.circular(20),
            border: Border.all(color: selected ? chipColor : context.appBorder),
          ),
          child: Text(
            label,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: selected ? Colors.white : context.appTextSecondary,
            ),
          ),
        ),
      ),
    );
  }
}

class _AssessmentHistoryCard extends StatelessWidget {
  final dynamic assessment;

  const _AssessmentHistoryCard({required this.assessment});

  @override
  Widget build(BuildContext context) {
    Color riskColor;
    switch (assessment.riskLevel) {
      case 'high':
        riskColor = AppColors.danger;
      case 'medium':
        riskColor = AppColors.warning;
      default:
        riskColor = AppColors.success;
    }

    final dateStr = DateFormat('MMM d, yyyy · h:mm a').format(assessment.date);

    return Card(
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: () => context.push('/patients/${assessment.patientId}'),
        child: Padding(
          padding: const EdgeInsets.all(14),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  CircleAvatar(
                    radius: 18,
                    backgroundColor: riskColor.withValues(alpha: 0.1),
                    child: Icon(
                      assessment.riskLevel == 'high'
                          ? Icons.warning_amber_rounded
                          : Icons.check_circle_outline_rounded,
                      color: riskColor,
                      size: 18,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          assessment.patientName,
                          style: const TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 17,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          dateStr,
                          style: Theme.of(context).textTheme.bodySmall,
                        ),
                      ],
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 3,
                    ),
                    decoration: BoxDecoration(
                      color: riskColor.withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      assessment.riskLevel.toUpperCase(),
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.w700,
                        color: riskColor,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 10),
              Row(
                children: [
                  Icon(
                    Icons.description_outlined,
                    size: 14,
                    color: context.appTextTertiary,
                  ),
                  const SizedBox(width: 6),
                  Text(
                    assessment.templateName,
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                  const Spacer(),
                  if (assessment.referralCreated) ...[
                    Icon(Icons.send_rounded, size: 12, color: AppColors.danger),
                    const SizedBox(width: 4),
                    Text(
                      'Referred',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: AppColors.danger,
                      ),
                    ),
                  ] else ...[
                    Icon(
                      Icons.cloud_upload_outlined,
                      size: 12,
                      color: context.appTextTertiary,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      'Pending sync',
                      style: TextStyle(
                        fontSize: 14,
                        color: context.appTextTertiary,
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
