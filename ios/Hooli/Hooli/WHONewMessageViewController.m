//
//  WHONewMessageViewController.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHONewMessageViewController.h"

@interface WHONewMessageViewController ()

@end

@implementation WHONewMessageViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // Custom initialization
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    
    UIBarButtonItem* submitButton = [[UIBarButtonItem alloc] initWithTitle:@"Done" style:UIBarButtonItemStylePlain target:self action:@selector(submit:)];
    self.navigationItem.rightBarButtonItem = submitButton;
    self.messageField.delegate = self;
    
    //for Done button on keyboard
//    [self.messageField addTarget:self action:@selector(messageView:) forControlEvents:UIControlEventEditingDidEndOnExit];
}

- (void)submit:(id)sender {
    //TODO compute user location
    [self.delegate receivedNewMessage:self.messageField.text];
    [self.navigationController popViewControllerAnimated:YES];
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

- (IBAction)messageView:(id)sender {
    [self.messageField resignFirstResponder];
}

-(void)textViewDidBeginEditing:(UITextView *)textView {
    self.placeholderLabel.hidden = YES;
}

//-(void)textViewDidChange:(UITextView *)textView {
//    self.placeholderLabel.hidden = [self.messageField.text length] > 0;
//}

-(void)textViewDidEndEditing:(UITextView *)textView {
    self.placeholderLabel.hidden = [self.messageField.text length] > 0;
}

//for Done button
- (BOOL)textView:(UITextView *)textView shouldChangeTextInRange:(NSRange)range replacementText:(NSString *)text {
    if([text isEqualToString:@"\n"]) {
        [textView resignFirstResponder];
        return NO;
    }
    return YES;
}

@end
