import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../core/theme/app_theme.dart';
import '../../core/data/mock_data.dart';
import '../../core/models/models.dart';

class AssessmentWizardScreen extends StatefulWidget {
  final String patientId;

  const AssessmentWizardScreen({super.key, required this.patientId});

  @override
  State<AssessmentWizardScreen> createState() => _AssessmentWizardScreenState();
}

class _AssessmentWizardScreenState extends State<AssessmentWizardScreen> {
  int _currentStep = 0;
  final int _totalSteps = 4;
  late Patient _patient;
  final AssessmentData _data = AssessmentData();

  // Controllers for vitals
  final _tempController = TextEditingController();
  final _weightController = TextEditingController();
  final _hrController = TextEditingController();
  final _spo2Controller = TextEditingController();
  final _notesController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _patient = MockData.patients.firstWhere(
      (p) => p.id == widget.patientId,
      orElse: () => MockData.patients.first,
    );
  }

  @override
  void dispose() {
    _tempController.dispose();
    _weightController.dispose();
    _hrController.dispose();
    _spo2Controller.dispose();
    _notesController.dispose();
    super.dispose();
  }

  void _next() {
    if (_currentStep == 0) {
      _data.temperature = double.tryParse(_tempController.text);
      _data.weight = double.tryParse(_weightController.text);
      _data.heartRate = int.tryParse(_hrController.text);
      _data.spo2 = int.tryParse(_spo2Controller.text);
    }
    if (_currentStep == 2) {
      _data.notes = _notesController.text;
    }

    if (_currentStep < _totalSteps - 1) {
      setState(() => _currentStep++);
    } else {
      // Navigate to result
      context.push(
        '/assessment-result',
        extra: {'assessmentData': _data, 'patientId': widget.patientId},
      );
    }
  }

  void _back() {
    if (_currentStep > 0) {
      setState(() => _currentStep--);
    } else {
      context.pop();
    }
  }

  @override
  Widget build(BuildContext context) {
    final stepTitles = ['Vitals', 'Symptoms', 'History & Notes', 'Review'];
    final stepSubtitles = [
      "Record the patient's vital signs",
      'Select all symptoms the patient reports',
      'Record relevant medical history',
      'Review all information before submitting',
    ];

    return Scaffold(
      backgroundColor: context.appBackground,
      appBar: AppBar(
        title: const Text('TB Screening'),
        leading: IconButton(
          icon: const Icon(Icons.close_rounded),
          onPressed: () => _showExitConfirmation(context),
        ),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16),
            child: Center(
              child: Text(
                'Step ${_currentStep + 1} of $_totalSteps',
                style: Theme.of(context).textTheme.bodySmall,
              ),
            ),
          ),
        ],
      ),
      body: Column(
        children: [
          // Progress bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(
                value: (_currentStep + 1) / _totalSteps,
                backgroundColor: context.appSurfaceDim,
                color: context.appPrimary,
                minHeight: 4,
              ),
            ),
          ),

          // Step header
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 20, 20, 8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  stepTitles[_currentStep],
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 4),
                Text(
                  stepSubtitles[_currentStep],
                  style: Theme.of(context).textTheme.bodyMedium,
                ),
              ],
            ),
          ),

          const SizedBox(height: 8),

          // Step content
          Expanded(
            child: AnimatedSwitcher(
              duration: const Duration(milliseconds: 200),
              child: _buildStepContent(),
            ),
          ),

          // Bottom buttons
          Container(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 32),
            decoration: BoxDecoration(
              color: context.appSurface,
              border: Border(top: BorderSide(color: context.appBorder)),
            ),
            child: Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: _back,
                    child: Text(_currentStep == 0 ? 'Cancel' : 'Back'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  flex: 2,
                  child: ElevatedButton(
                    onPressed: _next,
                    child: Text(
                      _currentStep == _totalSteps - 1 ? 'Submit' : 'Next',
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStepContent() {
    switch (_currentStep) {
      case 0:
        return _VitalsStep(
          key: const ValueKey('vitals'),
          tempController: _tempController,
          weightController: _weightController,
          hrController: _hrController,
          spo2Controller: _spo2Controller,
        );
      case 1:
        return _SymptomsStep(
          key: const ValueKey('symptoms'),
          selectedIds: _data.selectedSymptomIds,
          onChanged: (ids) => setState(() => _data.selectedSymptomIds = ids),
        );
      case 2:
        return _HistoryStep(
          key: const ValueKey('history'),
          data: _data,
          notesController: _notesController,
          onChanged: () => setState(() {}),
        );
      case 3:
        return _ReviewStep(
          key: const ValueKey('review'),
          data: _data,
          patient: _patient,
        );
      default:
        return const SizedBox.shrink();
    }
  }

  void _showExitConfirmation(BuildContext context) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Discard assessment?'),
        content: const Text('Your progress will be lost if you exit now.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Keep Editing'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              context.pop();
            },
            style: TextButton.styleFrom(foregroundColor: AppColors.danger),
            child: const Text('Discard'),
          ),
        ],
      ),
    );
  }
}

// ─── Step 1: Vitals ──────────────────────────────────────────

class _VitalsStep extends StatelessWidget {
  final TextEditingController tempController;
  final TextEditingController weightController;
  final TextEditingController hrController;
  final TextEditingController spo2Controller;

  const _VitalsStep({
    super.key,
    required this.tempController,
    required this.weightController,
    required this.hrController,
    required this.spo2Controller,
  });

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        children: [
          _VitalInput(
            icon: Icons.thermostat_rounded,
            iconColor: AppColors.danger,
            label: 'Temperature',
            suffix: '°C',
            controller: tempController,
            hint: '36.5',
          ),
          const SizedBox(height: 12),
          _VitalInput(
            icon: Icons.monitor_weight_rounded,
            iconColor: context.appPrimary,
            label: 'Weight',
            suffix: 'kg',
            controller: weightController,
            hint: '60',
          ),
          const SizedBox(height: 12),
          _VitalInput(
            icon: Icons.favorite_rounded,
            iconColor: AppColors.danger,
            label: 'Heart Rate',
            suffix: 'bpm',
            controller: hrController,
            hint: '72',
          ),
          const SizedBox(height: 12),
          _VitalInput(
            icon: Icons.air_rounded,
            iconColor: AppColors.info,
            label: 'SpO2',
            suffix: '%',
            controller: spo2Controller,
            hint: '98',
          ),
        ],
      ),
    );
  }
}

class _VitalInput extends StatelessWidget {
  final IconData icon;
  final Color iconColor;
  final String label;
  final String suffix;
  final TextEditingController controller;
  final String hint;

  const _VitalInput({
    required this.icon,
    required this.iconColor,
    required this.label,
    required this.suffix,
    required this.controller,
    required this.hint,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: context.appSurface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: context.appBorder),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(
              color: iconColor.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: iconColor, size: 22),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  label,
                  style: TextStyle(
                    fontSize: 14,
                    color: context.appTextSecondary,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                TextField(
                  controller: controller,
                  keyboardType: const TextInputType.numberWithOptions(
                    decimal: true,
                  ),
                  style: const TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w700,
                  ),
                  decoration: InputDecoration(
                    hintText: hint,
                    filled: false,
                    border: InputBorder.none,
                    enabledBorder: InputBorder.none,
                    focusedBorder: InputBorder.none,
                    contentPadding: EdgeInsets.zero,
                    isDense: true,
                  ),
                ),
              ],
            ),
          ),
          Text(
            suffix,
            style: TextStyle(
              fontSize: 18,
              color: context.appTextTertiary,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}

// ─── Step 2: Symptoms ────────────────────────────────────────

class _SymptomsStep extends StatelessWidget {
  final List<String> selectedIds;
  final ValueChanged<List<String>> onChanged;

  const _SymptomsStep({
    super.key,
    required this.selectedIds,
    required this.onChanged,
  });

  IconData _getIcon(String iconName) {
    switch (iconName) {
      case 'lungs':
        return Icons.air_rounded;
      case 'thermometer':
        return Icons.thermostat_rounded;
      case 'moon':
        return Icons.nightlight_round;
      case 'trending_down':
        return Icons.trending_down_rounded;
      case 'heart':
        return Icons.favorite_rounded;
      case 'battery_low':
        return Icons.battery_1_bar_rounded;
      case 'droplet':
        return Icons.water_drop_rounded;
      default:
        return Icons.circle;
    }
  }

  @override
  Widget build(BuildContext context) {
    final symptoms = MockData.tbSymptoms;

    return ListView.separated(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      itemCount: symptoms.length,
      separatorBuilder: (_, __) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final symptom = symptoms[index];
        final isSelected = selectedIds.contains(symptom.id);

        return GestureDetector(
          onTap: () {
            final ids = List<String>.from(selectedIds);
            if (isSelected) {
              ids.remove(symptom.id);
            } else {
              ids.add(symptom.id);
            }
            onChanged(ids);
          },
          child: AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: isSelected ? context.appPrimaryLight : context.appSurface,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(
                color: isSelected ? context.appPrimary : context.appBorder,
                width: isSelected ? 1.5 : 1,
              ),
            ),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: isSelected
                        ? context.appPrimary.withValues(alpha: 0.15)
                        : context.appSurfaceDim,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(
                    _getIcon(symptom.icon),
                    size: 20,
                    color: isSelected
                        ? context.appPrimary
                        : context.appTextSecondary,
                  ),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Text(
                    symptom.label,
                    style: TextStyle(
                      fontSize: 17,
                      fontWeight: isSelected
                          ? FontWeight.w600
                          : FontWeight.w400,
                      color: isSelected
                          ? context.appPrimary
                          : context.appTextPrimary,
                    ),
                  ),
                ),
                AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  width: 24,
                  height: 24,
                  decoration: BoxDecoration(
                    color: isSelected ? context.appPrimary : Colors.transparent,
                    borderRadius: BorderRadius.circular(6),
                    border: Border.all(
                      color: isSelected
                          ? context.appPrimary
                          : context.appBorder,
                      width: 1.5,
                    ),
                  ),
                  child: isSelected
                      ? const Icon(
                          Icons.check_rounded,
                          size: 16,
                          color: Colors.white,
                        )
                      : null,
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ─── Step 3: History & Notes ─────────────────────────────────

class _HistoryStep extends StatelessWidget {
  final AssessmentData data;
  final TextEditingController notesController;
  final VoidCallback onChanged;

  const _HistoryStep({
    super.key,
    required this.data,
    required this.notesController,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _ToggleCard(
            label: 'Has the patient been in contact with a TB case?',
            value: data.tbContactHistory,
            onChanged: (v) {
              data.tbContactHistory = v;
              onChanged();
            },
          ),
          const SizedBox(height: 12),
          _ToggleCard(
            label: 'Is the patient HIV positive?',
            value: data.hivPositive,
            onChanged: (v) {
              data.hivPositive = v;
              onChanged();
            },
          ),
          const SizedBox(height: 12),
          _ToggleCard(
            label: 'Has the patient been previously treated for TB?',
            value: data.previousTbTreatment,
            onChanged: (v) {
              data.previousTbTreatment = v;
              onChanged();
            },
          ),
          const SizedBox(height: 20),
          Text(
            'Clinical Notes',
            style: Theme.of(context).textTheme.titleMedium,
          ),
          const SizedBox(height: 8),
          TextField(
            controller: notesController,
            maxLines: 5,
            decoration: const InputDecoration(
              hintText: 'Add any additional observations...',
            ),
          ),
        ],
      ),
    );
  }
}

class _ToggleCard extends StatelessWidget {
  final String label;
  final bool? value;
  final ValueChanged<bool> onChanged;

  const _ToggleCard({
    required this.label,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      decoration: BoxDecoration(
        color: context.appSurface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: context.appBorder),
      ),
      child: Row(
        children: [
          Expanded(
            child: Text(
              label,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
            ),
          ),
          const SizedBox(width: 12),
          Row(
            children: [
              _OptionButton(
                label: 'Yes',
                isSelected: value == true,
                isPositive: true,
                onTap: () => onChanged(true),
              ),
              const SizedBox(width: 6),
              _OptionButton(
                label: 'No',
                isSelected: value == false,
                isPositive: false,
                onTap: () => onChanged(false),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _OptionButton extends StatelessWidget {
  final String label;
  final bool isSelected;
  final bool isPositive;
  final VoidCallback onTap;

  const _OptionButton({
    required this.label,
    required this.isSelected,
    required this.isPositive,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final color = isSelected
        ? (isPositive ? AppColors.danger : AppColors.success)
        : context.appTextTertiary;
    final bgColor = isSelected
        ? (isPositive ? context.appDangerLight : context.appSuccessLight)
        : Colors.transparent;

    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
        decoration: BoxDecoration(
          color: bgColor,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: isSelected ? color : context.appBorder),
        ),
        child: Text(
          label,
          style: TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w600,
            color: isSelected ? color : context.appTextTertiary,
          ),
        ),
      ),
    );
  }
}

// ─── Step 4: Review ──────────────────────────────────────────

class _ReviewStep extends StatelessWidget {
  final AssessmentData data;
  final Patient patient;

  const _ReviewStep({super.key, required this.data, required this.patient});

  @override
  Widget build(BuildContext context) {
    final symptoms = MockData.tbSymptoms
        .where((s) => data.selectedSymptomIds.contains(s.id))
        .toList();

    Color riskColor;
    switch (data.riskLevel) {
      case 'high':
        riskColor = AppColors.danger;
      case 'medium':
        riskColor = AppColors.warning;
      default:
        riskColor = AppColors.success;
    }

    return SingleChildScrollView(
      padding: const EdgeInsets.symmetric(horizontal: 20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Risk indicator
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: riskColor.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: riskColor.withValues(alpha: 0.3)),
            ),
            child: Row(
              children: [
                Icon(
                  data.shouldRefer
                      ? Icons.warning_amber_rounded
                      : Icons.check_circle_rounded,
                  color: riskColor,
                  size: 28,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Risk Level: ${data.riskLevel.toUpperCase()}',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w700,
                          color: riskColor,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        data.shouldRefer
                            ? 'Referral recommended'
                            : 'No referral needed',
                        style: TextStyle(
                          fontSize: 15,
                          color: riskColor.withValues(alpha: 0.8),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          const SizedBox(height: 20),

          // Patient
          _ReviewSection(label: 'Patient', value: patient.fullName),
          _ReviewSection(label: 'Village', value: patient.village ?? '—'),

          const Divider(height: 24),

          // Vitals
          Text('Vitals', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          _ReviewSection(
            label: 'Temperature',
            value: data.temperature != null
                ? '${data.temperature}°C'
                : 'Not recorded',
          ),
          _ReviewSection(
            label: 'Weight',
            value: data.weight != null ? '${data.weight} kg' : 'Not recorded',
          ),
          _ReviewSection(
            label: 'Heart Rate',
            value: data.heartRate != null
                ? '${data.heartRate} bpm'
                : 'Not recorded',
          ),
          _ReviewSection(
            label: 'SpO2',
            value: data.spo2 != null ? '${data.spo2}%' : 'Not recorded',
          ),

          const Divider(height: 24),

          // Symptoms
          Text(
            'Symptoms (${symptoms.length})',
            style: Theme.of(context).textTheme.titleMedium,
          ),
          const SizedBox(height: 8),
          if (symptoms.isEmpty)
            Text(
              'No symptoms reported',
              style: Theme.of(context).textTheme.bodySmall,
            ),
          ...symptoms.map(
            (s) => Padding(
              padding: const EdgeInsets.symmetric(vertical: 3),
              child: Row(
                children: [
                  Icon(
                    Icons.check_circle_rounded,
                    size: 16,
                    color: AppColors.danger,
                  ),
                  const SizedBox(width: 8),
                  Text(s.label, style: const TextStyle(fontSize: 16)),
                ],
              ),
            ),
          ),

          const Divider(height: 24),

          // History
          Text('History', style: Theme.of(context).textTheme.titleMedium),
          const SizedBox(height: 8),
          _ReviewSection(
            label: 'TB contact',
            value: data.tbContactHistory == true ? 'Yes' : 'No',
          ),
          _ReviewSection(
            label: 'HIV positive',
            value: data.hivPositive == true ? 'Yes' : 'No',
          ),
          _ReviewSection(
            label: 'Previous TB treatment',
            value: data.previousTbTreatment == true ? 'Yes' : 'No',
          ),
          if (data.notes.isNotEmpty) ...[
            const SizedBox(height: 8),
            Text('Notes', style: Theme.of(context).textTheme.bodySmall),
            const SizedBox(height: 4),
            Text(data.notes, style: const TextStyle(fontSize: 16)),
          ],

          const SizedBox(height: 40),
        ],
      ),
    );
  }
}

class _ReviewSection extends StatelessWidget {
  final String label;
  final String value;

  const _ReviewSection({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          SizedBox(
            width: 130,
            child: Text(label, style: Theme.of(context).textTheme.bodySmall),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
            ),
          ),
        ],
      ),
    );
  }
}
