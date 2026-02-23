import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';
import '../../core/models/models.dart';

class AssessmentResultScreen extends StatelessWidget {
  final AssessmentData assessmentData;
  final String patientId;

  const AssessmentResultScreen({
    super.key,
    required this.assessmentData,
    required this.patientId,
  });

  @override
  Widget build(BuildContext context) {
    final patient = MockData.patients.firstWhere(
      (p) => p.id == patientId,
      orElse: () => MockData.patients.first,
    );
    final isRefer = assessmentData.shouldRefer;
    final riskLevel = assessmentData.riskLevel;
    final criteria = assessmentData.matchedCriteria;

    Color riskColor;
    switch (riskLevel) {
      case 'high':
        riskColor = AppColors.danger;
      case 'medium':
        riskColor = AppColors.warning;
      default:
        riskColor = AppColors.success;
    }

    return Scaffold(
      backgroundColor: AppColors.background,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: Column(
            children: [
              const SizedBox(height: 20),

              // Success icon
              Container(
                width: 64,
                height: 64,
                decoration: BoxDecoration(
                  color: AppColors.successLight,
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.check_circle_rounded,
                  color: AppColors.success,
                  size: 36,
                ),
              ),
              const SizedBox(height: 16),
              Text(
                'Assessment Complete',
                style: Theme.of(context).textTheme.headlineMedium,
              ),

              const SizedBox(height: 24),

              // Result card
              if (isRefer)
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      colors: [
                        AppColors.danger,
                        AppColors.danger.withValues(alpha: 0.85),
                      ],
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                    ),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const Icon(
                            Icons.warning_amber_rounded,
                            color: Colors.white,
                            size: 24,
                          ),
                          const SizedBox(width: 10),
                          Text(
                            'Refer to Clinic',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 20,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 12),
                      Text(
                        'This patient shows ${criteria.length} red-flag indicators',
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.9),
                          fontSize: 14,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Text(
                        'Criteria met:',
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.8),
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        criteria.join(', '),
                        style: TextStyle(
                          color: Colors.white.withValues(alpha: 0.9),
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      const SizedBox(height: 12),
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                          vertical: 6,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.white.withValues(alpha: 0.2),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          'Template: Standard TB Screening v2.1',
                          style: TextStyle(
                            color: Colors.white.withValues(alpha: 0.9),
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ],
                  ),
                )
              else
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: riskColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: riskColor.withValues(alpha: 0.3),
                    ),
                  ),
                  child: Column(
                    children: [
                      Icon(
                        Icons.check_circle_outline_rounded,
                        color: riskColor,
                        size: 40,
                      ),
                      const SizedBox(height: 12),
                      Text(
                        riskLevel == 'medium'
                            ? 'Monitor Patient'
                            : 'Low Risk',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.w700,
                          color: riskColor,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        riskLevel == 'medium'
                            ? 'Schedule a follow-up in 2 weeks'
                            : 'No referral needed at this time',
                        style: TextStyle(
                          fontSize: 14,
                          color: riskColor.withValues(alpha: 0.8),
                        ),
                      ),
                    ],
                  ),
                ),

              const SizedBox(height: 24),

              // Summary details
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: AppColors.surface,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: AppColors.border),
                ),
                child: Column(
                  children: [
                    _SummaryRow(label: 'Patient', value: patient.fullName),
                    const Divider(height: 16),
                    _SummaryRow(
                      label: 'Assessment date',
                      value: '${DateTime.now().day}/${DateTime.now().month}/${DateTime.now().year}',
                    ),
                    const Divider(height: 16),
                    _SummaryRow(label: 'CHW', value: 'Sarah Kimani'),
                    const Divider(height: 16),
                    _SummaryRow(
                      label: 'Risk score',
                      value: riskLevel.toUpperCase(),
                      valueColor: riskColor,
                    ),
                    const Divider(height: 16),
                    _SummaryRow(
                      label: 'Recommendation',
                      value: isRefer
                          ? 'Refer to nearest clinic'
                          : 'No referral needed',
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 32),

              // Action buttons
              if (isRefer)
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton.icon(
                    onPressed: () {
                      _showReferralCreated(context);
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors.danger,
                      padding: const EdgeInsets.symmetric(vertical: 16),
                    ),
                    icon: const Icon(Icons.send_rounded),
                    label: const Text('Create Referral'),
                  ),
                ),
              const SizedBox(height: 10),
              SizedBox(
                width: double.infinity,
                child: OutlinedButton(
                  onPressed: () {
                    // Navigate back to home
                    context.go('/home');
                  },
                  child: Text(
                    isRefer ? 'Save Without Referral' : 'Done',
                  ),
                ),
              ),

              const SizedBox(height: 16),

              // Sync notice
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.cloud_off_rounded,
                    size: 14,
                    color: AppColors.textTertiary,
                  ),
                  const SizedBox(width: 6),
                  Text(
                    'This assessment will be synced when connected',
                    style: TextStyle(
                      fontSize: 12,
                      color: AppColors.textTertiary,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }

  void _showReferralCreated(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        icon: const Icon(
          Icons.check_circle_rounded,
          color: AppColors.success,
          size: 48,
        ),
        title: const Text('Referral Created'),
        content: const Text(
          'The referral has been saved and will be synced to the clinic when you have internet connectivity.',
        ),
        actions: [
          ElevatedButton(
            onPressed: () {
              Navigator.pop(ctx);
              context.go('/home');
            },
            child: const Text('Back to Home'),
          ),
        ],
      ),
    );
  }
}

class _SummaryRow extends StatelessWidget {
  final String label;
  final String value;
  final Color? valueColor;

  const _SummaryRow({
    required this.label,
    required this.value,
    this.valueColor,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: Text(
            label,
            style: Theme.of(context).textTheme.bodySmall,
          ),
        ),
        valueColor != null
            ? Container(
                padding:
                    const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                decoration: BoxDecoration(
                  color: valueColor!.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  value,
                  style: TextStyle(
                    fontSize: 13,
                    fontWeight: FontWeight.w700,
                    color: valueColor,
                  ),
                ),
              )
            : Text(
                value,
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
      ],
    );
  }
}
